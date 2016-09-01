package store

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

var (
	docKeyPrefix = []byte{'d', 'o', 'c', ':'}
)

func ignoreKey(key []byte) bool {
	//Keys starting with doc: are used internally, however, certain ip addresses like
	//100.111.99.58 and 646f:633a:... encode to 'doc:' in hex
	//Ignore the key if it starts with doc: unless it is exactly 4 or 16 bytes long
	if bytes.HasPrefix(key, docKeyPrefix) {
		kl := len(key)
		if kl != 4 && kl != 16 {
			return true
		}
	}
	return false
}

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
	Filename() string
}

var DefaultStore = "leveldb"

var storeFactories = map[string]func(string) (IpStore, error){}

func NewStore(storeType string, filename string) (IpStore, error) {
	storeFactory, ok := storeFactories[storeType]
	if !ok {
		return nil, errors.New("Invalid store type")
	}
	s, err := storeFactory(filename)
	return s, err
}
