#!/bin/bash
# Example of compsing commands to dump everything in the same format
# Useful for packing the full history
# Could be used in map reduce scenarios
# $1 is the path to the rocksdb directory
# $2 is the start hash
# $3 is the end hash
# example invocation:
# ./dump.sh ~/ripple/nodedb/ E6DB7365949BF9814D76BCC730B01818EB9136A89DB224F3F9F5AAE4569D758E E6DB7365949BF9814D76BCC730B01818EB9136A89DB224F3F9F5AAE4569D758E | gzip -c9 >dump.gz
rdb -path=$1 -command=ledgers -start=$2 -end=$3
rdb -path=$1 -command=transactions -start=$2 -end=$3 -dump_format="%[4]s,%[5]X"
rdb -path=$1 -command=diff -start=$2 -end=$3 -diff_format="%[4]s,%[5]X"
