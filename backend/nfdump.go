package backend

import (
	"io"
	"os/exec"

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
	lines, err := SyslogBackend{}.ExtractIps(stdout, ips)
	cmd.Wait()
	return lines, err

}

func init() {
	RegisterBackend("nfdump", NFDUMPBackend{})
}
