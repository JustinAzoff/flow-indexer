package backend

import (
	"io"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

const maxLineLength = 20 * 1024 * 1024

type Backend interface {
	ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error)
	Filter(reader io.Reader, query string, writer io.Writer) error
	Check() error
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
	reader, err := OpenDecompress(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ExtractIpsReader(backend, reader)
}

func ExtractIpsReader(backend string, reader io.Reader) (*ipset.Set, error) {
	ips := ipset.New()
	b := NewBackend(backend)
	_, err := b.ExtractIps(reader, ips)
	return ips, err
}

func FilterIPs(backend string, filename string, query string, writer io.Writer) error {
	reader, err := OpenDecompress(filename)
	if err != nil {
		return err
	}
	defer reader.Close()
	return FilterIPsReader(backend, reader, query, writer)
}

func FilterIPsReader(backend string, reader io.Reader, query string, writer io.Writer) error {
	b := NewBackend(backend)
	err := b.Filter(reader, query, writer)
	return err
}
