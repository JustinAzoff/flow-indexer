package backend

import (
	"testing"
)

var expectedBroOutput = []string{}
var expectedBroCount = 12

func TBroExtractIps(backend string, t *testing.T) {
	ips, err := ExtractIps(backend, "test_data/bro_conn.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	if ips.Count() != expectedBroCount {
		t.Errorf("BroBackend(%s).ExtractIps count => %#v, want %#v", backend, ips.Count(), expectedBroCount)
	}
	for _, ip := range ips.SortedStrings() {
		t.Logf("%x\n", ip)
	}
}
func TestBroNativeExtractIps(t *testing.T) {
	TBroExtractIps("bronative", t)
}
func TestBroCutExtractIps(t *testing.T) {
	TBroExtractIps("brocut", t)
}

func BenchBroExtract(backend string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractIps(backend, "test_data/bro_conn_some_v6.log.gz")
	}
}
func BenchmarkBroNativeExtract(b *testing.B) {
	BenchBroExtract("bronative", b)
}
func BenchmarkBroCutExtract(b *testing.B) {
	BenchBroExtract("brocut", b)
}
