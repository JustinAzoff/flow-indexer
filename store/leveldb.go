package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/willf/bitset"
)

type LevelDBStore struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

func NewLevelDBStore(filename string) (IpStore, error) {
	db, err := leveldb.OpenFile(filename, nil)
	if err != nil {
		return nil, err
	}
	newStore := &LevelDBStore{db: db, batch: nil}
	return newStore, nil
}

func (ls *LevelDBStore) Close() error {
	return ls.db.Close()
}

func (ls *LevelDBStore) HasDocument(filename string) (bool, error) {
	_, err := ls.db.Get([]byte(filename), nil)
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
	ls.setDocId(filename, nextID)
	ls.batch = new(leveldb.Batch)
	for k, _ := range ips.Store {
		//fmt.Printf("Add %#v to document\n", k)
		ls.addIP(nextID, k)
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
			break
		}
		fmt.Printf("Document %d is %#v\n", i, name)
	}
	return err
}

func (ls *LevelDBStore) DocumentIDToName(id uint64) (string, error) {
	idBytes := PutUVarint(id)
	v, err := ls.db.Get(idBytes, nil)
	return string(v), err
}

func (ls *LevelDBStore) ExpandCIDR(ip string) ([]net.IP, error) {
	var ips []net.IP
	start, end, err := ipset.CIDRToByteStrings(ip)
	if err != nil {
		return ips, err
	}
	iter := ls.db.NewIterator(&util.Range{Start: []byte(start), Limit: []byte(end)}, nil)
	for iter.Next() {
		key := iter.Key()
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
	bs := bitset.New(8)
	tmpbs := bitset.New(8)
	iter := ls.db.NewIterator(&util.Range{Start: []byte(start), Limit: []byte(end)}, nil)
	for iter.Next() {
		tmpbs.ReadFrom(bytes.NewBuffer(iter.Value()))
		bs = bs.Union(tmpbs)
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
	key, err := ipset.IPToByteString(ip)
	if err != nil {
		return nil, err
	}
	v, err := ls.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return docs, nil
	}
	bs := bitset.New(8)
	bs.ReadFrom(bytes.NewBuffer(v))
	return ls.bitsetToDocs(bs)
}

func (ls *LevelDBStore) bitsetToDocs(bs *bitset.BitSet) ([]string, error) {
	var err error
	var docs []string
	for i, e := bs.NextSet(0); e; i, e = bs.NextSet(i + 1) {
		name, err := ls.DocumentIDToName(uint64(i))
		if err != nil {
			break
		}
		docs = append(docs, name)
	}
	return docs, err
}

func (ls *LevelDBStore) nextDocID() (uint64, error) {
	v, err := ls.db.Get([]byte("max_id"), nil)
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
func (ls *LevelDBStore) setDocId(filename string, id uint64) error {
	idBytes := PutUVarint(id)
	ls.db.Put([]byte(filename), idBytes, nil)
	ls.db.Put(idBytes, []byte(filename), nil)
	return ls.db.Put([]byte("max_id"), idBytes, nil)
}

func (ls *LevelDBStore) addIP(id uint64, k string) error {
	bs := bitset.New(8)
	v, err := ls.db.Get([]byte(k), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	if err != leveldb.ErrNotFound {
		bs.ReadFrom(bytes.NewBuffer(v))
	}
	bs.Set(uint(id))

	buffer := bytes.NewBuffer(make([]byte, 0, bs.BinaryStorageSize()))
	_, err = bs.WriteTo(buffer)
	if err != nil {
		return err
	}
	ls.batch.Put([]byte(k), buffer.Bytes())
	return nil
}
