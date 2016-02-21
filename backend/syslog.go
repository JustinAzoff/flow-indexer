package backend

import (
	"bufio"
	"io"
	"log"
	"regexp"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type SyslogBackend struct {
}

var IPRegexString = `(?:[^0-9](?P<ip>[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})[^0-9])|` +
	`(?:[^0-9A-Fa-f](?P<ip>(([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4}))[^0-9A-Fa-f])|` +
	`(?:[^0-9A-Fa-f](?P<ip>(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?)::(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?))[^0-9A-Fa-f])|` + //IPv6 Compressed Hex
	`(?:[^0-9A-Fa-f](?P<ip>(([0-9A-Fa-f]{1,4}:){6,6})([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+))[^0-9A-Fa-f])|` + //6Hex4Dec
	`(?:[^0-9A-Fa-f](?P<ip>(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?)::(([0-9A-Fa-f]{1,4}:)*)([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+))[^0-9A-Fa-f])` //CompressedHex4Dec

var IPRegex = regexp.MustCompile(IPRegexString)

var IPIndexes []int

func init() {
	for i, name := range IPRegex.SubexpNames() {
		if name == "ip" {
			IPIndexes = append(IPIndexes, i)
		}
	}
}

func (b SyslogBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	br := bufio.NewReader(reader)

	lines := uint64(0)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return lines, err
		}
		lines++
		ipsFound := IPRegex.FindAllStringSubmatch(line, -1)

		for _, ipMatches := range ipsFound {
			for _, idx := range IPIndexes {
				if ipMatches[idx] != "" {
					ips.AddString(ipMatches[idx])
				}
			}
		}
	}
	return lines, nil
}

func init() {
	RegisterBackend("syslog", SyslogBackend{})
}
