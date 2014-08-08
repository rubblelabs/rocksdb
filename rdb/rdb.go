package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/ledger"
	"github.com/rubblelabs/ripple/storage"
	"github.com/rubblelabs/ripple/storage/memdb"
	"github.com/rubblelabs/rocksdb"
)

var (
	unfortunateGenesis = mustDecodeHash("4109C6F2045FC7EFF4CDE8F9905D19C28820D86304080FF886B299F0206E42B5") // 32,570
	defaultStart       = mustDecodeHash("491E88B0A5AB29378B4F4E6EAB1E782AF495D712A817C943D0D7A36045EFA611") // 100,000
)

var start = flag.String("start", defaultStart.String(), "initial ledger hash (defaults to 100000)")
var end = flag.String("end", unfortunateGenesis.String(), "final ledger hash (defaults to 32570)")
var path = flag.String("path", "", "location of db folder")
var mem = flag.Bool("mem", false, "use memory db")
var cpu = flag.Int("cpu", runtime.NumCPU()*2, "number of cpu's to use")
var command = flag.String("command", "diff", "command to run [diff/dump/summary/transaction/ledgers/accounts]")
var cacheSize = flag.Int("cache_size", 0, "size of the rocksdb memory cache")
var dump_format = flag.String("dump_format", "%[1]d,%[2]s,%[3]d,%[4]s,%[5]X", "customisable format string for dump/accounts/transactions commands. Indexes: [1]:LedgerSequence, [2]:NodeType, [3]:Depth, [4]:NodeId, [5]:NodeValue")
var diff_format = flag.String("diff_format", "%[1]d,%[6]c,%[2]s,%[3]d,%[4]s,%[5]X", "customisable format string for diff command Indexes: [1]:LedgerSequence, [2]:NodeType, [3]:Depth, [4]:NodeId, [5]:NodeValue, [6]:Action")

func checkErr(err error) {
	if err != nil {
		glog.Errorln(err.Error())
		os.Exit(1)
	}
}

var db storage.DB

var commands = map[string]func() ledger.QueueFunc{
	"diff":         diff,
	"dump":         dump,
	"summary":      summary,
	"ledgers":      ledgers,
	"transactions": transactions,
	"accounts":     accounts,
}

func nodeFields(n *ledger.RadixNode, ledgerSequence uint32) ([]interface{}, error) {
	nodeId, nodeValue, err := data.Node(n.Node)
	if err != nil {
		return nil, err
	}
	return []interface{}{ledgerSequence, n.Node.GetType(), n.Depth, nodeId, nodeValue}, nil
}

func outputNode(n *ledger.RadixNode, ledgerSequence uint32) error {
	fields, err := nodeFields(n, ledgerSequence)
	if err != nil {
		return err
	}
	_, err = fmt.Printf(*dump_format+"\n", fields...)
	return err
}

func outputDiffNode(op *ledger.RadixOperation, ledgerSequence uint32) error {
	fields, err := nodeFields(op.RadixNode, ledgerSequence)
	if err != nil {
		return err
	}
	_, err = fmt.Printf(*diff_format+"\n", append(fields, op.Action)...)
	return err
}

func outputLedger(h data.Storer) error {
	hash, value, err := data.Node(h)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s,%X\n", hash, value)
	return err
}

func dumpAccounts(includeInner bool) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		if err := current.AccountState.Fill(); err != nil {
			return err
		}
		return current.AccountState.Walk(func(key data.Hash256, node *ledger.RadixNode) error {
			if _, ok := node.Node.(data.LedgerEntry); ok || includeInner {
				return outputNode(node, current.LedgerSequence)
			}
			return nil
		})
	}
}

func dumpTransactions(includeInner bool) ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		if err := current.Transactions.Fill(); err != nil {
			// Temp fix for bad memos
			_, err := fmt.Fprintf(os.Stderr, "Skipping ledger: %d %s\n", current.LedgerSequence, err.Error())
			return err
		}
		return current.Transactions.Walk(func(key data.Hash256, node *ledger.RadixNode) error {
			txm, ok := node.Node.(*data.TransactionWithMetaData)
			if ok {
				// Fix for bad ledger sequences in nodestore
				txm.LedgerSequence = current.LedgerSequence
			}
			if ok || includeInner {
				return outputNode(node, current.LedgerSequence)
			}
			return nil
		})
	}
}

func accounts() ledger.QueueFunc {
	return dumpAccounts(false)
}

func transactions() ledger.QueueFunc {
	return dumpTransactions(false)
}

func dump() ledger.QueueFunc {
	txDumper := dumpTransactions(true)
	accountDumper := dumpAccounts(true)
	return func(current, previous *ledger.LedgerState) error {
		current.Fill()
		if err := txDumper(current, previous); err != nil {
			return err
		}
		return accountDumper(current, previous)
	}
}

func ledgers() ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		return outputLedger(current.Ledger)
	}
}

func diff() ledger.QueueFunc {
	return func(current, previous *ledger.LedgerState) error {
		diff, err := ledger.Diff(current.StateHash, previous.StateHash, db)
		if err != nil {
			return err
		}
		for _, op := range diff {
			if err := outputDiffNode(op, current.LedgerSequence); err != nil {
				return err
			}
		}
		return nil
	}
}

func summary() ledger.QueueFunc {
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

func mustDecodeHash(hash string) *data.Hash256 {
	h, err := data.NewHash256(hash)
	checkErr(err)
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
	cmd, ok := commands[*command]
	if !ok {
		glog.Errorln("Unknown command:", *command)
		os.Exit(1)
	}
	var queue ledger.Queue
	for {
		state, err := ledger.NewLedgerStateFromDB(*from, db)
		checkErr(err)
		queue.Add(state)
		checkErr(queue.Do(cmd()))
		if *from == *to {
			if *command != "diff" || *to == *unfortunateGenesis {
				queue.AddEmpty()
				checkErr(queue.Do(cmd()))
			}
			return
		}
		*from = state.PreviousLedger
	}
}
