package backend

import "testing"

var expectedPCAPOutput = []string{}
var expectedPCAPCount = 23

func TestPCAPExtractIps(t *testing.T) {
	ips, err := ExtractIps("pcap", "test_data/pcap.pcap.gz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedPCAPCount {
		t.Errorf("PCAPBackend.ExtractIps count => %#v, want %#v", ips.Count(), expectedPCAPCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}

	t.Logf("%v\n", ips.Count())
}

var testRandomASCIIPCAPLogBytes []byte

func BenchmarkPCAPExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("pcap", "test_data/pcap.pcap.gz")
	}
}
