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

var IPRegexString = `[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}|` +
	`([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4}|` +
	`(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?)::(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?)|` + //IPv6 Compressed Hex
	`(([0-9A-Fa-f]{1,4}:){6,6})([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+)|` + //6Hex4Dec
	`(([0-9A-Fa-f]{1,4}(:[0-9A-Fa-f]{1,4})*)?)::(([0-9A-Fa-f]{1,4}:)*)([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+)` //CompressedHex4Dec
var IPRegex = regexp.MustCompile(IPRegexString)

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
		ipsFound := IPRegex.FindAllString(line, -1)
		for _, ip := range ipsFound {
			ips.AddString(ip)
		}
	}
	return lines, nil
}

func init() {
	RegisterBackend("syslog", SyslogBackend{})
}
