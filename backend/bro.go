package backend

import (
	"bufio"
	"io"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type BroBackend struct {
}

func (b BroBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	br := bufio.NewReader(reader)

	lines := uint64(0)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		if line[0] != '#' {
			parts := strings.SplitN(line, "\t", 6) //makes parts[4] the last full split
			ips.AddString(parts[2])
			ips.AddString(parts[4])
			lines++
		}
	}
	return lines, nil
}

func init() {
	RegisterBackend("bro", BroBackend{})
}
