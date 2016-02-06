package main

import (
	"log"
	"time"

	"github.com/justinazoff/flow-indexer/backend"
	"github.com/justinazoff/flow-indexer/store"
)

func Index(s store.IpStore, b backend.Backend, filename string) error {
	exists, err := s.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("%s Already indexed\n", filename)
		return nil
	}

	ips, err := b.ExtractIps(filename)
	if err != nil {
		return err
	}
	start := time.Now()
	s.AddDocument(filename, *ips)
	duration := time.Since(start)
	log.Printf("%s: Wrote %d unique ips in %s\n", filename, len(ips.Store), duration)
	return nil
}
