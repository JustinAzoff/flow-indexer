package backend

import (
	"bytes"
	"testing"
)

var expectedSyslogOutput = []string{}
var expectedSyslogCount = 12

func TestSyslogExtractIps(t *testing.T) {
	ips, err := ExtractIps("syslog", "test_data/bro_conn_some_v6.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedBroCount {
		t.Errorf("SyslogBackend.ExtractIps count => %#v, want %#v", ips.Count(), expectedBroCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}

	t.Logf("%v\n", ips.Count())
}

func BenchmarkSyslogExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps("syslog", "test_data/bro_conn_some_v6.log.gz")
	}
}

func BenchmarkSyslogExtractRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIpsReader("bro", bytes.NewBuffer(testRandomASCIIBroLogBytes))
	}
}
