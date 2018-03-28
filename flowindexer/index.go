package flowindexer

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/JustinAzoff/flow-indexer/backend"
	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/JustinAzoff/flow-indexer/store"
)

func Index(s store.IpStore, b backend.Backend, filename string) error {
	exists, err := s.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		//log.Printf("%s Already indexed\n", filename)
		return nil
	}

	if err = b.Check(); err != nil {
		return fmt.Errorf("Backend is not usable: %v", err)
	}

	ips := ipset.New()
	reader, err := backend.OpenDecompress(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	start := time.Now()
	lines, err := b.ExtractIps(reader, ips)
	duration := time.Since(start)
	log.Printf("%s: Read %d lines in %s\n", filename, lines, duration)

	if err != nil {
		log.Printf("%s: Non fatal read error: %s\n", filename, err)
	}
	start = time.Now()
	err = s.AddDocument(filename, *ips)
	if err != nil {
		return err
	}
	duration = time.Since(start)
	log.Printf("%s: Wrote %d unique ips in %s\n", filename, ips.Count(), duration)
	return nil
}

func RunIndex(dbpath string, backend_type string, args []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mystore, err := store.NewStore("leveldb", dbpath)
	check(err)
	defer mystore.Close()

	for _, path := range args {
		mybackend := backend.NewBackend(backend_type)
		matches, err := filepath.Glob(path)
		check(err)
		for _, fp := range matches {
			err = Index(mystore, mybackend, fp)
			check(err)
		}
	}
}
