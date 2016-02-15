package backend

import (
	"testing"
)

var expectedBroOutput = []string{}
var expectedBroCount = 12

func TestBroExtractIps(t *testing.T) {
	ips, err := ExtractIps("bro", "test_data/bro_conn.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	if len(ips.Store) != expectedBroCount {
		t.Errorf("BroBackend.ExtractIps count => %#v, want %#v", len(ips.Store), expectedBroCount)
	}
	for k, _ := range ips.Store {
		t.Logf("%x\n", k)
	}

	t.Logf("%v\n", len(ips.Store))
}

func BenchmarkBroExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("bro", "test_data/bro_conn_some_v6.log.gz")
	}
}
