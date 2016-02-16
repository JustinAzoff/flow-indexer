//This converts the bitset based values to msgpack encoded slices
//What I found was that the raw size of the values was much smaller:
//
//bitset size is 1129122792
//msgpack size is 497539014
//
//But the post-snappy compressed size did not change
//
//428M    2015-03.db
//433M    2015-03.test
//
//The bitsets are mostly zeroes, so this is not unexpected.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/willf/bitset"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func Open(filename string) (*leveldb.DB, error) {
	fmt.Printf("Opening %s\n", filename)
	//Options taken from ledisdb
	opts := &opt.Options{}
	opts.BlockSize = 32768
	opts.WriteBuffer = 67108864
	opts.BlockCacheCapacity = 524288000
	opts.OpenFilesCacheCapacity = 1024
	opts.CompactionTableSize = 32 * 1024 * 1024
	opts.WriteL0SlowdownTrigger = 16
	opts.WriteL0PauseTrigger = 64
	opts.Filter = filter.NewBloomFilter(10)

	db, err := leveldb.OpenFile(filename, opts)
	return db, err
}

func main() {

	olddb, err := Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	newdb, err := Open(os.Args[2])
	if err != nil {
		panic(err)
	}

	iter := olddb.NewIterator(&util.Range{Start: nil, Limit: nil}, nil)

	totalBitset := 0
	totalMsgpack := 0

	rows := 0

	var batch *leveldb.Batch
	batch = new(leveldb.Batch)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if bytes.HasPrefix(key, []byte("doc:")) {
			batch.Put(key, value)
			continue
		}
		bs := bitset.New(8)
		bs.ReadFrom(bytes.NewBuffer(value))
		var docIDs []uint
		for i, e := bs.NextSet(0); e; i, e = bs.NextSet(i + 1) {
			docIDs = append(docIDs, i)
		}
		b, err := msgpack.Marshal(docIDs)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("bitset size is %d\n", len(value))
		//fmt.Printf("msgpack size is %d\n", len(b))

		totalBitset += len(value)
		totalMsgpack += len(b)
		batch.Put(key, b)
		if rows%10000 == 0 {
			log.Print("rows ", rows)
			newdb.Write(batch, nil)
			batch = new(leveldb.Batch)
		}
		rows++

	}
	fmt.Printf("bitset size is %d\n", totalBitset)
	fmt.Printf("msgpack size is %d\n", totalMsgpack)
	newdb.Write(batch, nil)
	newdb.CompactRange(util.Range{Start: nil, Limit: nil})

}
