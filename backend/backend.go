package backend

import (
	"github.com/JustinAzoff/flow-indexer/ipset"
)

type ExtractResult struct {
	records int64
	set     ipset.Set
}

type Backend interface {
	ExtractIps(filename string) (*ipset.Set, error)
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
