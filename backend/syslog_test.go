package backend

import (
	"testing"
)

var expectedSyslogOutput = []string{}
var expectedSyslogCount = 12

func TestSyslogExtractIps(t *testing.T) {
	b := SyslogBackend{}
	ips, err := b.ExtractIps("test_data/bro_conn_some_v6.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	if len(ips.Store) != expectedBroCount {
		t.Errorf("SyslogBackend.ExtractIps count => %#v, want %#v", len(ips.Store), expectedBroCount)
	}
	for k, _ := range ips.Store {
		t.Logf("%x\n", k)
	}

	t.Logf("%v\n", len(ips.Store))
}
