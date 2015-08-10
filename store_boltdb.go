package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/justinazoff/flow-indexer/ipset"
	"github.com/willf/bitset"
)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(filename string) (*BoltStore, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.CreateBucketIfNotExists([]byte("docs"))
	if err != nil {
		return nil, fmt.Errorf("create bucket: %s", err)
	}
	_, err = tx.CreateBucketIfNotExists([]byte("ips"))
	if err != nil {
		return nil, fmt.Errorf("create bucket: %s", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	newStore := &BoltStore{db}
	return newStore, nil

}

func (bs *BoltStore) Close() error {
	return bs.db.Close()
}

func (bs *BoltStore) HasDocument(filename string) (bool, error) {
	var exists bool
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("docs"))
		v := b.Get([]byte(filename))
		exists = v != nil
		return nil
	})
	if err != nil {
		return exists, err
	}
	return exists, nil
}

func (bs *BoltStore) AddDocument(filename string, ips ipset.Set) error {
	exists, err := bs.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("docs"))
		nextID, err := nextDocID(b)
		if err != nil {
			return err
		}
		fmt.Printf("NextDocID should be %d\n", nextID)
		setDocId(b, filename, nextID)
		ipBucket := tx.Bucket([]byte("ips"))
		for k, _ := range ips.Store {
			//fmt.Printf("Add %#v to document\n", k)
			addIP(ipBucket, nextID, k)
		}
		return nil
	})

	return nil
}

func (bs *BoltStore) ListDocuments() error {
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("docs"))
		nextID, err := nextDocID(b)
		for i := uint64(0); i < nextID; i += 1 {
			name := DocumentIDToName(b, i)
			fmt.Printf("Document %d is %#v\n", i, name)
		}
		return err
	})
	return err
}

func DocumentIDToName(b *bolt.Bucket, id uint64) string {
	idBytes := PutUVarint(id)
	v := b.Get(idBytes)
	return string(v)
}

func (bs *BoltStore) QueryString(ip string) {
	key, err := ipset.IPToByteString(ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("docs"))
		ipBucket := tx.Bucket([]byte("ips"))
		v := ipBucket.Get([]byte(key))
		if v == nil {
			fmt.Printf("%s does not exist\n", ip)
			return nil
		}
		bs := bitset.New(8)
		bs.ReadFrom(bytes.NewBuffer(v))
		for i, e := bs.NextSet(0); e; i, e = bs.NextSet(i + 1) {
			fmt.Printf("Match in document %d %s\n", i, DocumentIDToName(b, uint64(i)))
		}
		return nil
	})
	return
}

func nextDocID(b *bolt.Bucket) (uint64, error) {
	v := b.Get([]byte("max_id"))
	if v == nil {
		return 0, nil
	}
	maxID, read := binary.Uvarint(v)
	if read <= 0 {
		return 0, fmt.Errorf("Error converting %#v to a uint64", v)
	}
	return maxID + 1, nil

}
func setDocId(b *bolt.Bucket, filename string, id uint64) error {
	idBytes := PutUVarint(id)
	b.Put([]byte(filename), idBytes)
	b.Put(idBytes, []byte(filename))
	return b.Put([]byte("max_id"), idBytes)
}

func addIP(b *bolt.Bucket, id uint64, k string) {
	v := b.Get([]byte(k))
	bs := bitset.New(8)
	if v != nil {
		bs.ReadFrom(bytes.NewBuffer(v))
	}
	bs.Set(uint(id))

	buffer := bytes.NewBuffer(make([]byte, 0, bs.BinaryStorageSize()))
	_, err := bs.WriteTo(buffer)
	if err != nil {
		return //nil, err
	}
	b.Put([]byte(k), buffer.Bytes())
}

func PutUVarint(v uint64) []byte {
	b := make([]byte, 10)
	binary.PutUvarint(b, uint64(v))
	return b
}
