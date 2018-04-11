package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/willf/bitset"
)

type LevelDBStore struct {
	filename     string
	db           *leveldb.DB
	batch        *leveldb.Batch
	codecFactory func() Codec
}

func NewLevelDBStore(filename string) (IpStore, error) {
	//Options taken from ledisdb
	opts := &opt.Options{}
	opts.BlockSize = 32768
	opts.BlockCacheCapacity = 524288000
	opts.OpenFilesCacheCapacity = 1024
	opts.CompactionTableSize = 32 * 1024 * 1024
	opts.WriteL0SlowdownTrigger = 16
	opts.WriteL0PauseTrigger = 64
	opts.Filter = filter.NewBloomFilter(10)

	db, err := leveldb.OpenFile(filename, opts)
	if err != nil {
		return nil, err
	}
	newStore := &LevelDBStore{db: db, batch: nil, filename: filename, codecFactory: func() Codec { return NewBitsetCodec() }}
	newStore.fixDocId()
	return newStore, nil
}

func (ls *LevelDBStore) Close() error {
	return ls.db.Close()
}

func (ls *LevelDBStore) Compact() error {
	return ls.db.CompactRange(util.Range{Start: nil, Limit: nil})
}

func (ls *LevelDBStore) Filename() string {
	return ls.filename
}

func (ls *LevelDBStore) HasDocument(filename string) (bool, error) {
	key := buildFilenameKey(filename)
	_, err := ls.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ls *LevelDBStore) AddDocument(filename string, ips ipset.Set) error {
	exists, err := ls.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	nextID, err := ls.nextDocID()
	if err != nil {
		return err
	}
	ls.batch = new(leveldb.Batch)
	ls.setDocId(filename, nextID)
	for _, ip := range ips.SortedStrings() {
		//fmt.Printf("Add %#v to document\n", ip)
		err = ls.addIP(nextID, ip)
		if err != nil {
			return err
		}
	}
	err = ls.db.Write(ls.batch, nil)
	ls.batch = nil
	return err

}

func (ls *LevelDBStore) ListDocuments() error {
	nextID, err := ls.nextDocID()
	for i := uint64(0); i < nextID; i += 1 {
		name, err := ls.DocumentIDToName(i)
		if err != nil {
			return err
		}
		fmt.Printf("Document %d is %#v\n", i, name)
	}
	return err
}

func (ls *LevelDBStore) DocumentIDToName(id uint64) (string, error) {
	idBytes := buildDocumentKey(id)
	v, err := ls.db.Get(idBytes, nil)
	return string(v), err
}

func (ls *LevelDBStore) ExpandCIDR(ip string) ([]net.IP, error) {
	var ips []net.IP
	start, end, err := ipset.CIDRToByteStrings(ip)
	if err != nil {
		return ips, err
	}
	bstart := []byte(start)
	bend := []byte(end)
	iter := ls.db.NewIterator(&util.Range{Start: bstart, Limit: nil}, nil)
	for iter.Next() {
		key := iter.Key()
		if ignoreKey(key) {
			continue
		}
		if len(key) != len(start) {
			//Ensure the matched keys are in the right ip family
			continue
		}
		if bytes.Compare(key, bend) > 0 {
			break
		}
		keycopy := make([]byte, len(key))
		copy(keycopy, key)
		ip := net.IP(keycopy)
		ips = append(ips, ip)
	}
	iter.Release()
	err = iter.Error()
	return ips, err
}

func (ls *LevelDBStore) QueryString(ip string) ([]string, error) {
	if strings.Contains(ip, "/") {
		return ls.QueryStringCidr(ip)
	}
	return ls.QueryStringIP(ip)
}

func (ls *LevelDBStore) QueryStringCidr(ip string) ([]string, error) {
	var start, end string
	start, end, err := ipset.CIDRToByteStrings(ip)
	if err != nil {
		return nil, err
	}
	bstart := []byte(start)
	bend := []byte(end)
	bs := bitset.New(8)
	iter := ls.db.NewIterator(&util.Range{Start: bstart, Limit: nil}, nil)

	codec := ls.codecFactory()
	for iter.Next() {
		key := iter.Key()
		if ignoreKey(key) {
			continue
		}
		if len(key) != len(start) {
			//Ensure the matched keys are in the right ip family
			continue
		}
		if bytes.Compare(key, bend) > 0 {
			break
		}
		codec.FromBytes(iter.Value())
		tmpBs := codec.ToBitset()
		bs = bs.Union(tmpBs)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		return nil, err
	}

	return ls.bitsetToDocs(bs)
}

func (ls *LevelDBStore) QueryStringIP(ip string) ([]string, error) {
	var docs []string
	key, err := ipset.IPStringToByteString(ip)
	if err != nil {
		return nil, err
	}
	v, err := ls.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return docs, nil
	}
	codec := ls.codecFactory()
	codec.FromBytes(v)
	bs := codec.ToBitset()
	return ls.bitsetToDocs(bs)
}

func (ls *LevelDBStore) bitsetToDocs(bs *bitset.BitSet) ([]string, error) {
	var docs []string
	for i, e := bs.NextSet(0); e; i, e = bs.NextSet(i + 1) {
		name, err := ls.DocumentIDToName(uint64(i))
		if err != nil {
			return docs, err
		}
		docs = append(docs, name)
	}
	return docs, nil
}

//fixDocId fixes an issue where the max docid was stored under max_id
//instead of doc:max_id so a search for 109.97.120.95 would find it
func (ls *LevelDBStore) fixDocId() {
	v, err := ls.db.Get([]byte("max_id"), nil)
	if err == leveldb.ErrNotFound {
		return
	}
	if err != nil {
		return
	}
	//key max_id exists, rewrite it to doc:max_id
	log.Println("FIX: Renaming max_id to doc:max_id")
	batch := new(leveldb.Batch)
	batch.Put([]byte("doc:max_id"), v)
	batch.Delete([]byte("max_id"))
	ls.db.Write(batch, nil)
}

func (ls *LevelDBStore) nextDocID() (uint64, error) {
	v, err := ls.db.Get([]byte("doc:max_id"), nil)
	if err == leveldb.ErrNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	maxID, read := binary.Uvarint(v)
	if read <= 0 {
		return 0, fmt.Errorf("Error converting %#v to a uint64", v)
	}
	return maxID + 1, nil

}

func (ls *LevelDBStore) setDocId(filename string, id uint64) {
	idBytes := buildDocumentKey(id) // doc:xxx
	fnBytes := buildFilenameKey(filename)
	ls.batch.Put(fnBytes, idBytes[4:])
	ls.batch.Put(idBytes, []byte(filename))
	ls.batch.Put([]byte("doc:max_id"), idBytes[4:])
}

func (ls *LevelDBStore) addIP(id uint64, k string) error {
	v, err := ls.db.Get([]byte(k), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	codec := ls.codecFactory()
	if err != leveldb.ErrNotFound {
		codec.FromBytes(v)
	}
	codec.AddID(DocumentID(id))

	bytes, err := codec.Bytes()
	if err != nil {
		return err
	}
	ls.batch.Put([]byte(k), bytes)
	return nil
}

func init() {
	storeFactories["leveldb"] = NewLevelDBStore
}
