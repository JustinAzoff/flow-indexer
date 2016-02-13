package backend

import (
	"testing"
)

var expectedBroOutput = []string{}
var expectedBroCount = 12

func TestBroExtractIps(t *testing.T) {
	b := BroBackend{}
	ips, err := b.ExtractIps("test_data/bro_conn.log.gz")
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
