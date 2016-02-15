package backend

import (
	"io"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type Backend interface {
	ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error)
}

var backends = map[string]Backend{}

func RegisterBackend(name string, backend Backend) {
	backends[name] = backend
}

func NewBackend(backendType string) Backend {
	backend, ok := backends[backendType]
	if !ok {
		panic("Invalid backend")
	}
	return backend
}

func ExtractIps(backend string, filename string) (*ipset.Set, error) {
	ips := ipset.New()
	reader, err := OpenDecompress(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	b := NewBackend(backend)
	_, err = b.ExtractIps(reader, ips)
	return ips, err
}
