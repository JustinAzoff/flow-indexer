package backend

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/JustinAzoff/flow-indexer/ipset"
)

type NFDUMPCSVBackend struct {
}

func (b NFDUMPCSVBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	cmd := exec.Command("nfdump", "-qr", "-", "-o", "csv")
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
			if err = ips.AddString(parts[3]); err != nil {
				return lines, err
			}
			if err = ips.AddString(parts[4]); err != nil {
				return lines, err
			}
			lines++
		}
	}
	err = cmd.Wait()
	return lines, err
}
func (b NFDUMPCSVBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	filter := fmt.Sprintf("ip in [%s]", query)
	cmd := exec.Command("nfdump", "-qr", "-", filter)
	cmd.Stdin = reader
	cmd.Stdout = writer

	err := cmd.Run()
	return err
}

func (b NFDUMPCSVBackend) Check() error {
	cmd := exec.Command("nfdump", "-V")
	_, err := cmd.CombinedOutput()
	return err
}

type NFDUMPBackend struct {
}

func (b NFDUMPBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	cmd := exec.Command("nfdump", "-qr", "-", "-o", "fmt:%sa %da")
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
func (b NFDUMPBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	filter := fmt.Sprintf("ip in [%s]", query)
	cmd := exec.Command("nfdump", "-qr", "-", filter)
	cmd.Stdin = reader
	cmd.Stdout = writer

	err := cmd.Run()
	return err
}

func (b NFDUMPBackend) Check() error {
	cmd := exec.Command("nfdump", "-V")
	_, err := cmd.CombinedOutput()
	return err
}

func init() {
	RegisterBackend("nfdump-csv", NFDUMPCSVBackend{})
	RegisterBackend("nfdump", NFDUMPBackend{})
}
