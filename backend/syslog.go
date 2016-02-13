package backend

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"time"

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

func (b SyslogBackend) ExtractIps(filename string) (*ipset.Set, error) {
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
		lines++
		ips := IPRegex.FindAllString(line, -1)
		for _, ip := range ips {
			s.AddString(ip)
		}
	}
	duration := time.Since(start)
	log.Printf("%s: Read %d lines in %s\n", filename, lines, duration)
	return s, nil
}

func init() {
	RegisterBackend("syslog", SyslogBackend{})
}
