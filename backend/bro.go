package backend

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/justinazoff/flow-indexer/ipset"
	gzip "github.com/klauspost/pgzip"
)

type BroBackend struct {
}

func (b BroBackend) ExtractIps(filename string) (*ipset.Set, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
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
			log.Print(err)
			return s, err
		}
		if line[0] != '#' {
			parts := strings.Split(line, "\t")
			s.AddString(parts[2])
			s.AddString(parts[4])
			lines++
		}
	}
	duration := time.Since(start)
	log.Printf("%s: Read %d lines in %s\n", filename, lines, duration)
	return s, nil
}
