//go:generate ffjson $GOFILE

package backend

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

//ffjson: skip
type BroJSONBackend struct {
}

type BroIPFields struct {
	ID_orig_h string `json:"id.orig_h"`
	ID_resp_h string `json:"id.resp_h"`
	Src       string `json:"src"`
	Dst       string `json:"dst"`
}

func (b BroJSONBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	br := bufio.NewReaderSize(reader, maxLineLength)

	lines := uint64(0)
	for {
		var FoundIPS BroIPFields
		line, err := br.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}

		err = FoundIPS.UnmarshalJSON(line)
		if err != nil {
			return lines, err
		}
		if FoundIPS.ID_orig_h != "" {
			ips.AddString(FoundIPS.ID_orig_h)
		}
		if FoundIPS.ID_resp_h != "" {
			ips.AddString(FoundIPS.ID_resp_h)
		}
		if FoundIPS.Src != "" {
			ips.AddString(FoundIPS.Src)
		}
		if FoundIPS.Dst != "" {
			ips.AddString(FoundIPS.Dst)
		}
		lines++
	}
	return lines, nil
}

func (b BroJSONBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	br := bufio.NewReaderSize(reader, maxLineLength)

	realQuery := fmt.Sprintf("\"%s\"", query)

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
func (b BroJSONBackend) Check() error {
	return nil
}

func init() {
	RegisterBackend("bro_json", BroJSONBackend{})
}
