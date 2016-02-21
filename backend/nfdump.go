package backend

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type NFDUMPBackend struct {
}

func (b NFDUMPBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	cmd := exec.Command("nfdump", "-r", "-", "-o", "csv")
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
		parts := strings.SplitN(line, ",", 6) //makes parts[4] the last full split
		if len(parts) == 6 {
			ips.AddString(parts[3])
			ips.AddString(parts[4])
			lines++
		}
	}
	err = cmd.Wait()
	return lines, err
}

func init() {
	RegisterBackend("nfdump", NFDUMPBackend{})
}
