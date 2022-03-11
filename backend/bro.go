package backend

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type BroBackend struct {
}

func (b BroBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	br := bufio.NewReaderSize(reader, maxLineLength)

	lines := uint64(0)
	var fields []string
	var ipFields []int
	var highestFieldIndex int
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		if line[0] == '#' && strings.HasPrefix(line, "#fields") {
			fields = strings.Split(line, "\t")
			fields = fields[1:]
			ipFields = ipFields[:0]
			for i, f := range fields {
				if f == "id.orig_h" || f == "id.resp_h" {
					ipFields = append(ipFields, i)
					highestFieldIndex = i
				}
			}
		}
		if line[0] != '#' {
			parts := strings.SplitN(line, "\t", highestFieldIndex+2) //Split just enough fields
			for _, idx := range ipFields {
				ips.AddString(parts[idx])
			}
			lines++
		}
	}
	return lines, nil
}

func (b BroBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	br := bufio.NewReaderSize(reader, maxLineLength)

	realQuery := fmt.Sprintf("\t%s\t", query)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if strings.Index(line, realQuery) != -1 {
			if _, err = io.WriteString(writer, line); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b BroBackend) Check() error {
	return nil
}

func init() {
	RegisterBackend("bro", BroBackend{})
}
