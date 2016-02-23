package backend

import (
	"bytes"
	"testing"

	"github.com/JustinAzoff/flow-indexer/loggen"
)

var expectedBroJSONOutput = []string{}
var expectedBroJSONCount = 12

func TestBroJSONExtractIps(t *testing.T) {
	ips, err := ExtractIps("bro_json", "test_data/bro_conn.log.json.gz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedBroJSONCount {
		t.Errorf("BroJSONBackend.ExtractIps count => %#v, want %#v", ips.Count(), expectedBroCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}

	t.Logf("%v\n", ips.Count())
}

var testRandomJSONBroLogBytes []byte

func init() {
	testRandomJSONBroLogBytes = loggen.RandomJSONBroLog(1000)
}

func BenchmarkBroJSONExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("bro_json", "test_data/bro_conn.log.json.gz")
	}
}

func BenchmarkBroJSONExtractRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIpsReader("bro", bytes.NewBuffer(testRandomJSONBroLogBytes))
	}
}
