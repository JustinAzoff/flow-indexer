package backend

import (
	"bytes"
	"testing"

	"github.com/JustinAzoff/flow-indexer/loggen"
)

var expectedBroOutput = []string{}
var expectedBroCount = 12

func TestBroExtractIps(t *testing.T) {
	ips, err := ExtractIps("bro", "test_data/bro_conn.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedBroCount {
		t.Errorf("BroBackend.ExtractIps count => %#v, want %#v", ips.Count(), expectedBroCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}

	t.Logf("%v\n", ips.Count())
}

var testRandomASCIIBroLogBytes []byte

func init() {
	testRandomASCIIBroLogBytes = loggen.RandomASCIIBroLog(1000)
}

func BenchmarkBroExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIpsReader("bro", bytes.NewBuffer(testRandomASCIIBroLogBytes))
	}
}

func BenchmarkBroExtractRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("bro", "test_data/bro_conn.log.gz")
	}
}
