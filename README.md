rocksdb
=======

Storage layer to access rippled nodestore and a tool for dumping its contents.


##Installation

Install [Go](http://golang.org/doc/install) making sure to set a GOPATH.
Install [RocksDB dependencies](https://github.com/facebook/rocksdb/blob/master/INSTALL.md).

```bash
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
make shared_lib
go get -u -v github.com/rubblelabs/rocksdb
```

##Usage

For OSX

```bash
export DYLD_LIBRARY_PATH=/path/to/rocksdb/
rdb -help
```

For Linux

```bash
export LD_LIBRARY_PATH=/path/to/rocksdb/
rdb -help
```
