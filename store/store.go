package store

import (
	"encoding/binary"
	"errors"

	"github.com/justinazoff/flow-indexer/ipset"
)

func PutUVarint(v uint64) []byte {
	b := make([]byte, 10)
	binary.PutUvarint(b, uint64(v))
	return b
}

type IpStore interface {
	Close() error
	HasDocument(filename string) (bool, error)
	AddDocument(filename string, ips ipset.Set) error
	QueryString(ip string) error
}

var storeFactories = map[string]func(string) (IpStore, error){
	"leveldb": NewLevelDBStore,
	//	"boltdb":  NewBoltStore,
}

func NewStore(storeType string, filename string) (IpStore, error) {
	storeFactory, ok := storeFactories[storeType]
	if !ok {
		return nil, errors.New("Invalid store type")
	}
	s, err := storeFactory(filename)
	return s, err
}
