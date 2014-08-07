package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/ledger"
	"github.com/rubblelabs/ripple/storage"
	"github.com/rubblelabs/ripple/storage/memdb"
	"github.com/rubblelabs/rocksdb"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

var (
	unfortunateGenesis = mustDecodeHash("4109C6F2045FC7EFF4CDE8F9905D19C28820D86304080FF886B299F0206E42B5") //32,570
	defaultStart       = mustDecodeHash("491E88B0A5AB29378B4F4E6EAB1E782AF495D712A817C943D0D7A36045EFA611") // 100,000
)

var start = flag.String("start", defaultStart.String(), "initial ledger hash (defaults to 100000)")
var end = flag.String("end", unfortunateGenesis.String(), "final ledger hash (defaults to 32570)")
var path = flag.String("path", "", "location of db folder")
var mem = flag.Bool("mem", false, "use memory db")
var cpu = flag.Int("cpu", runtime.NumCPU()*2, "number of cpu's to use")
var command = flag.String("command", "diff", "command to run [diff/historypack/dump/summary/transaction/ledgers/accounts]")
var cacheSize = flag.Int("cache_size", 0, "size of the rocksdb memory cache")

func checkErr(err error) {
	if err != nil {
		glog.Errorln(err.Error())
		os.Exit(1)
	}
}

var db storage.DB

var commands = map[string]func(io.Writer) ledger.QueueFunc{
	"diff":         diff,
	"historypack":  historypack,
	"dump":         dump,
	"summary":      summary,
	"transactions": transactions,
	"ledgers":      ledgers,
	"accounts":     accounts,
}

func writeHex(w io.Writer, h data.Storer) error {
	hash, value, err := data.Node(h)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s:%X\n", hash, value)
	return err
}

func accounts(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		if err := current.AccountState.Fill(); err != nil {
			return err
		}
		return current.AccountState.Walk(func(key data.Hash256, node *ledger.RadixNode) error {
			if _, ok := node.Node.(data.LedgerEntry); ok {
				return writeHex(w, node.Node)
			}
			return nil
		})
	}
}

func transactions(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		if err := current.Transactions.Fill(); err != nil {
			// Temp fix for bad memos
			_, err := fmt.Fprintf(os.Stderr, "Skipping ledger: %d %s\n", current.LedgerSequence, err.Error())
			return err
		}
		return current.Transactions.Walk(func(key data.Hash256, node *ledger.RadixNode) error {
			if txm, ok := node.Node.(*data.TransactionWithMetaData); ok {
				// Fix for bad ledger sequences in nodestore
				txm.LedgerSequence = current.LedgerSequence
				return writeHex(w, txm)
			}
			return nil
		})
	}
}

func ledgers(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		return writeHex(w, current.Ledger)
	}
}

func diff(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		diff, err := ledger.Diff(current.StateHash, previous.StateHash, db)
		if err != nil {
			return err
		}
		for _, op := range diff {
			if _, err := fmt.Fprintf(w, "%d,%s\n", current.LedgerSequence, op); err != nil {
				return err
			}
		}
		return nil
	}
}

func historypack(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		diff, err := ledger.Diff(current.StateHash, previous.StateHash, db)
		if err != nil {
			return err
		}
		for _, op := range diff {
			if _, err := fmt.Fprintf(w, "%d:%c:", current.LedgerSequence, op.Action); err != nil {
				return err
			}
			if err := writeHex(w, op.Node); err != nil {
				return err
			}
		}
		return nil
	}
}

func dump(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		current.Fill()
		if err := current.Transactions.Dump(current.LedgerSequence, w); err != nil {
			return err
		}
		return current.AccountState.Dump(current.LedgerSequence, w)
	}
}

func summary(w io.Writer) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		current.Fill()
		summary, err := current.Summary()
		if err != nil {
			return err
		}
		_, err = fmt.Printf("%d,%s\n", current.LedgerSequence, summary)
		return err
	}
}

func do(from, to *data.Hash256, f ledger.QueueFunc) error {
	var queue ledger.Queue
	for {
		state, err := ledger.NewLedgerStateFromDB(*from, db)
		if err != nil {
			return err
		}
		queue.Add(state)
		if err := queue.Do(f); err != nil {
			return err
		}
		if from.Compare(*to) == 0 {
			if from.Compare(*unfortunateGenesis) == 0 {
				queue.AddEmpty()
			}
			return queue.Do(f)
		}
		*from = state.PreviousLedger
	}
}

func mustDecodeHash(hash string) *data.Hash256 {
	h, err := data.NewHash256(hash)
	if err != nil {
		glog.Fatalln(err.Error())
	}
	return h
}

func report() {
	thirty := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-thirty.C:
			glog.Infoln(db.Stats())
		}
	}
}

func main() {
	flag.Parse()
	from, to := mustDecodeHash(*start), mustDecodeHash(*end)
	runtime.GOMAXPROCS(*cpu)
	go func() {
		glog.Infoln(http.ListenAndServe("localhost:6060", nil))
	}()
	var err error
	if *mem {
		db, err = memdb.NewMemoryDB([]string{*path})
	} else {
		db, err = rocksdb.NewRocksDB(*path, *cacheSize)
	}
	checkErr(err)
	defer db.Close()
	go report()
	if cmd, ok := commands[*command]; ok {
		checkErr(do(from, to, cmd(os.Stdout)))
	} else {
		glog.Errorln("Unknown command: ", *command)
	}
}
