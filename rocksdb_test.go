package rocksdb

import (
	"testing"
)

// Just tests the thing compiles and returns an error
func TestRocksDB(t *testing.T) {
	db, err := NewRocksDB(".", 100)
	if db != nil {
		t.Errorf("Expected nil db")
	}
	if err == nil {
		t.Errorf("Expected err")
	}
}
