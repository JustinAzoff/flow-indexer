package backend

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/justinazoff/flow-indexer/ipset"
	"github.com/justinazoff/flow-indexer/store"
	gzip "github.com/klauspost/pgzip"
)

func IndexBroLog(store store.IpStore, filename string) error {
	exists, err := store.HasDocument(filename)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("%s Already indexed\n", filename)
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	br := bufio.NewReader(reader)

	s := ipset.New()
	lines := 0
	start := time.Now()
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print(err)
			return err
		}
		if line[0] != '#' {
			parts := strings.Split(line, "\t")
			s.AddString(parts[2])
			s.AddString(parts[4])
			lines++
		}
	}
	duration := time.Since(start)
	fmt.Printf("%s: Read %d lines in %s\n", filename, lines, duration)

	start = time.Now()
	store.AddDocument(filename, *s)
	duration = time.Since(start)
	fmt.Printf("%s: Wrote %d unique ips in %s\n", filename, len(s.Store), duration)
	return nil
}
