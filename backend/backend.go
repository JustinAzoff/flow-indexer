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

var backends = map[string]Backend{
	"bro": BroBackend{},
}

func NewBackend(backendType string) Backend {
	backend, ok := backends[backendType]
	if !ok {
		panic("Invalid backend")
	}
	return backend
}
