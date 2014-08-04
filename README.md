rocksdb
=======

Storage layer to access rippled nodestore and a tool for dumping its contents.


##Installation

Install [Go](http://golang.org/doc/install) making sure to set a GOPATH and add GOPATH/bin to your PATH.
Install [RocksDB dependencies](https://github.com/facebook/rocksdb/blob/master/INSTALL.md).

```bash
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
make shared_lib
CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb" go get -u -v github.com/rubblelabs/rocksdb/rdb
```

changing the /path/to/rocksdb as appropriate.


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
