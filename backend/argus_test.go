package backend

import "testing"

var expectedArgusOutput = []string{}
var expectedArgusCount = 8

func TestArgusExtractIps(t *testing.T) {
	ips, err := ExtractIps("argus", "test_data/argus.data.xz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedArgusCount {
		t.Errorf("ArgusBackend.ExtractIps count => %#v, want %#v", ips.Count(), expectedArgusCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}

	t.Logf("%v\n", ips.Count())
}

var testRandomASCIIArgusLogBytes []byte

func BenchmarkArgusExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("argus", "test_data/argus.data.xz")
	}
}
