package backend

import (
	"encoding/json"
	"io"
	"log"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type BroJSONBackend struct {
}
type BroIPFields struct {
	ID_orig_h string `json:"id.orig_h"`
	ID_resp_h string `json:"id.resp_h"`
	Src       string `json:"src"`
	Dst       string `json:"dst"`
}

func (b BroJSONBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	dec := json.NewDecoder(reader)
	lines := uint64(0)
	for {
		var FoundIPS BroIPFields
		err := dec.Decode(&FoundIPS)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
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

func init() {
	RegisterBackend("bro_json", BroJSONBackend{})
}
