package backend

import (
	"testing"
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

func BenchmarkBroJSONExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("bro_json", "test_data/bro_conn.log.json.gz")
	}
}
