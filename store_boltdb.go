package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/justinazoff/flow-indexer/ipset"
	"log"
)

type BoltStore struct {
	db         *bolt.DB
	docsBucket *bolt.Bucket
	ipsBucket  *bolt.Bucket
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
	docsBucket, err := tx.CreateBucketIfNotExists([]byte("docs"))
	if err != nil {
		return nil, fmt.Errorf("create bucket: %s", err)
	}
	ipsBucket, err := tx.CreateBucketIfNotExists([]byte("ips"))
	if err != nil {
		return nil, fmt.Errorf("create bucket: %s", err)
	}
	newStore := &BoltStore{
		db:         db,
		docsBucket: docsBucket,
		ipsBucket:  ipsBucket,
	}
	return newStore, nil
}

func (bs *BoltStore) Close() error {
	return bs.db.Close()
}

func (bs *BoltStore) AddDocument(filename string, ips ipset.Set) error {
	return nil
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	bs, err := NewBoltStore("my.db")
	check(err)
	defer bs.Close()
}
