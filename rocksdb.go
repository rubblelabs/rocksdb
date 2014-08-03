package rocksdb

import (
	"bytes"
	"fmt"
	"github.com/golang/groupcache/lru"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/storage"
	"github.com/tecbot/gorocksdb"
	"sync"
)

type RocksDB struct {
	db           *gorocksdb.DB
	ro           *gorocksdb.ReadOptions
	cache        *lru.Cache
	hits, misses metrics.Meter
	mu           sync.RWMutex
}

func NewRocksDB(path string, cacheSize int) (*RocksDB, error) {
	opts := gorocksdb.NewDefaultOptions()
	filter := gorocksdb.NewBloomFilter(14)
	opts.SetFilterPolicy(filter)
	opts.SetMaxOpenFiles(10000)
	db, err := gorocksdb.OpenDbForReadOnly(opts, path, false)
	if err != nil {
		return nil, err
	}
	return &RocksDB{
		db:     db,
		ro:     gorocksdb.NewDefaultReadOptions(),
		hits:   metrics.NewMeter(),
		misses: metrics.NewMeter(),
		cache:  lru.New(cacheSize),
	}, nil
}

func (db *RocksDB) Close() error {
	db.db.Close()
	return nil
}

func (db *RocksDB) Get(hash data.Hash256) (data.Storer, error) {
	db.mu.Lock()
	cached, ok := db.cache.Get(hash)
	db.mu.Unlock()
	if ok {
		db.hits.Mark(1)
		return cached.(data.Storer), nil
	}
	value, err := db.db.Get(db.ro, hash[:])
	if err != nil {
		return nil, err
	}
	defer value.Free()
	if value.Size() == 0 {
		return nil, storage.ErrNotFound
	}
	node, err := data.ReadPrefix(bytes.NewReader(value.Data()), hash)
	if err != nil {
		return nil, err
	}
	db.misses.Mark(1)
	db.mu.Lock()
	db.cache.Add(hash, node)
	db.mu.Unlock()
	return node, nil
}

func (db *RocksDB) Insert(data.Storer) error {
	panic("not implemented")
}

func (db *RocksDB) Ledger() (*data.LedgerSet, error) {
	return data.NewLedgerSet(32570, 32570), nil
}

func (db *RocksDB) Stats() string {
	db.mu.RLock()
	entries := db.cache.Len()
	size := db.cache.MaxEntries
	db.mu.RUnlock()
	hits, misses := db.hits.Snapshot(), db.misses.Snapshot()
	total := float64(hits.Count() + misses.Count())
	hitsPercent, missesPercent := float64(hits.Count())/total*100, float64(misses.Count())/total*100
	occupancy := float64(entries) / float64(size) * 100
	format := "Cache: %0.02f%% full Hits: %0.02f%% %d(%.0f/sec) Misses: %0.02f%%  %d(%.0f/sec)"
	return fmt.Sprintf(format, occupancy, hitsPercent, hits.Count(), hits.RateMean(), missesPercent, misses.Count(), misses.RateMean())
}
