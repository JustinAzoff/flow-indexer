package backend

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type PCAPBackend struct {
}

func (b PCAPBackend) ExtractIps(reader io.Reader, ips *ipset.Set) (uint64, error) {
	packets := uint64(0)

	pr, err := pcapgo.NewReader(reader)
	if err != nil {
		return 0, err
	}
	var eth layers.Ethernet
	var dot1q layers.Dot1Q
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &dot1q, &ip4, &ip6, &tcp)
	decoded := []gopacket.LayerType{}
	for {
		packetData, _, err := pr.ReadPacketData()
		packets++
		if err == io.EOF {
			break
		}
		if err != nil {
			return packets, err
		}
		err = parser.DecodeLayers(packetData, &decoded)
		for _, layerType := range decoded {
			switch layerType {
			case layers.LayerTypeIPv6:
				ips.AddIP(ip6.SrcIP)
				ips.AddIP(ip6.DstIP)
			case layers.LayerTypeIPv4:
				ips.AddIP(ip4.SrcIP)
				ips.AddIP(ip4.DstIP)
			}
		}
	}

	return packets, nil
}

func (b PCAPBackend) Filter(reader io.Reader, query string, writer io.Writer) error {
	filter := fmt.Sprintf("(net %s) or (vlan and net %s) or (vlan and vlan and net %s)", query, query, query)

	cmd := exec.Command("tcpdump", "-nn", "-r", "-", filter)
	cmd.Stdin = reader
	cmd.Stdout = writer

	err := cmd.Run()
	return err
}

func (b PCAPBackend) Check() error {
	return nil
}

func init() {
	RegisterBackend("pcap", PCAPBackend{})
}
