package backend

import (
	"bufio"
	"io"
	"log"
	"strings"
	"time"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type BroBackend struct {
}

func (b BroBackend) ExtractIps(filename string) (*ipset.Set, error) {
	reader, err := OpenDecompress(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
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
			parts := strings.SplitN(line, "\t", 6) //makes parts[4] the last full split
			s.AddString(parts[2])
			s.AddString(parts[4])
			lines++
		}
	}
	duration := time.Since(start)
	log.Printf("%s: Read %d lines in %s\n", filename, lines, duration)
	return s, nil
}

func init() {
	RegisterBackend("bro", BroBackend{})
}
