package backend

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

var broCutAvailable bool

func init() {
	cmd := exec.Command("bro-cut")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		return
	}
	broCutAvailable = true
	log.Printf("Using bro-cut for bro field extraction")
}

type BroBackend struct {
}

func (b BroBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	if broCutAvailable {
		return BroCutBackend{}.ExtractIps(reader, ips)
	} else {
		return BroNativeBackend{}.ExtractIps(reader, ips)
	}
}

type BroCutBackend struct {
}

func (b BroCutBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	cmd := exec.Command("bro-cut", "id.orig_h", "id.resp_h")
	cmd.Stdin = reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	err = cmd.Start()
	if err != nil {
		return 0, err
	}

	br := bufio.NewReader(stdout)

	lines := uint64(0)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		parts := strings.Split(line, "\t")
		for _, ip := range parts {
			ips.AddString(ip)
		}
		lines++
	}
	err = cmd.Wait()
	return lines, err
}

type BroNativeBackend struct {
}

func (b BroNativeBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {

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
	RegisterBackend("bronative", BroNativeBackend{})
	RegisterBackend("brocut", BroCutBackend{})
}
