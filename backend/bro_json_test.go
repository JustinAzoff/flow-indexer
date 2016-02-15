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
	if len(ips.Store) != expectedBroJSONCount {
		t.Errorf("BroJSONBackend.ExtractIps count => %#v, want %#v", len(ips.Store), expectedBroCount)
	}
	for k, _ := range ips.Store {
		t.Logf("%x\n", k)
	}

	t.Logf("%v\n", len(ips.Store))
}

func BenchmarkBroJSONExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("bro_json", "test_data/bro_conn.log.json.gz")
	}
}
