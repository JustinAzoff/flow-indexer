package backend

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type ArgusBackend struct {
}

func (b ArgusBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	cmd := exec.Command("ra", "-s", "saddr,daddr", "-c", " ", "-L", "-1")
	cmd.Stdin = reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	err = cmd.Start()
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)

	lines := uint64(0)
	for scanner.Scan() {
		if err = ips.AddString(scanner.Text()); err != nil {
			return lines / 2, err
		}
		lines++
	}
	if err := scanner.Err(); err != nil {
		return lines / 2, err
	}
	err = cmd.Wait()

	// Each line gets counted twice in the scanner.Scan for loop
	return lines / 2, err
}
func (b ArgusBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	filter := fmt.Sprintf("net %s", query)
	cmd := exec.Command("ra", "--", filter)
	cmd.Stdin = reader
	cmd.Stdout = writer

	err := cmd.Run()
	return err
}

func (b ArgusBackend) Check() error {
	//ra is annoying and always exits with status code 1.
	cmd := exec.Command("ra", "-h")
	out, err := cmd.CombinedOutput()
	if bytes.Contains(out, []byte("Ra Version")) {
		return nil
	}
	return err
}

func init() {
	RegisterBackend("argus", ArgusBackend{})
}
