package store

import (
	"encoding/binary"
	"errors"
	"github.com/JustinAzoff/flow-indexer/ipset"
	"net"
)

var (
	docKeyPrefix = []byte{'d', 'o', 'c', ':'}
)

func PutUVarint(v uint64) []byte {
	b := make([]byte, 10)
	binary.PutUvarint(b, uint64(v))
	return b
}

//buildDocumentKey builds a byte array containing doc: and the varint encoded id
func buildDocumentKey(id uint64) []byte {
	b := make([]byte, 10+len(docKeyPrefix))
	copy(b[:], docKeyPrefix)
	binary.PutUvarint(b[len(docKeyPrefix):], id)
	return b
}

func buildFilenameKey(fn string) []byte {
	b := make([]byte, len(docKeyPrefix)+len(fn))
	copy(b[:], docKeyPrefix)
	copy(b[len(docKeyPrefix):], fn)
	return b
}

type IpStore interface {
	Close() error
	HasDocument(filename string) (bool, error)
	AddDocument(filename string, ips ipset.Set) error
	QueryString(ip string) ([]string, error)
	ExpandCIDR(ip string) ([]net.IP, error)
	Compact() error
}

var storeFactories = map[string]func(string) (IpStore, error){
	"leveldb": NewLevelDBStore,
	//	"boltdb":  NewBoltStore,
}

func NewStore(storeType string, filename string) (IpStore, error) {
	_, ok := storeFactories[storeType]
	if !ok {
		return nil, errors.New("Invalid store type")
	}
	s, err := NewLevelDBStore(filename)
	return s, err
}
