package store

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/JustinAzoff/flow-indexer/loggen"
)

func makeIps(ss []string) []net.IP {
	ips := make([]net.IP, len(ss))
	for i, s := range ss {
		ips[i] = net.ParseIP(s)
	}

	return ips
}

var basicSearchTable = []struct {
	query string
	docs  []string
}{
	{"1.2.3.4/24", []string{"/log/1.txt", "/log/2.txt"}},
	{"2.0.0.0/8", []string{"/log/2.txt"}},
}

var basicExpandCidrTable = []struct {
	query string
	ips   []net.IP
}{
	{"1.2.3.0/24", makeIps([]string{"1.2.3.1", "1.2.3.2", "1.2.3.3", "1.2.3.4"})},
	{"1.0.0.0/8", makeIps([]string{"1.2.3.1", "1.2.3.2", "1.2.3.3", "1.2.3.4"})},
	{"2.0.0.0/8", makeIps([]string{"2.0.0.2", "2.0.0.3"})},

	//'doc:' converted to an IP is 100.111.99.58
	{"100.111.99.0/24", []net.IP{}},
	//'/log' converted to an IP is 47.108.111.103
	{"47.108.111.0/24", []net.IP{}},
	//max_id converted to an ip is 109.97.120.95
	{"109.97.120.0/24", []net.IP{}},
}

func runTest(t *testing.T, s IpStore) {
	ips := ipset.New()
	ips.AddString("1.2.3.1")
	ips.AddString("1.2.3.2")
	ips.AddString("1.2.3.3")
	ips.AddString("1.2.3.4")
	s.AddDocument("/log/1.txt", *ips)

	ips.AddString("2.0.0.2")
	ips.AddString("2.0.0.3")
	s.AddDocument("/log/2.txt", *ips)

	for _, tt := range basicSearchTable {
		matches, err := s.QueryString(tt.query)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(matches, tt.docs) {
			t.Errorf("store.QueryString => %#v, want %#v", matches, tt.docs)
		}
	}

	for _, tt := range basicExpandCidrTable {
		matches, err := s.ExpandCIDR(tt.query)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%v", matches) != fmt.Sprintf("%v", tt.ips) {
			t.Errorf("store.QueryString => %v, want %v", matches, tt.ips)
		}
	}

}

func TestLeveldb(t *testing.T) {
	mystore, err := NewStore("leveldb", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer mystore.Close()
	defer os.RemoveAll("test.db")
	runTest(t, mystore)

}

func runStoreBench(b *testing.B, storeType string) {
	mystore, err := NewStore(storeType, "test.db")
	if err != nil {
		b.Error(err)
		return
	}
	defer os.RemoveAll("test.db")
	for doc := 0; doc < 24; doc++ {
		ips := ipset.New()
		for i := 0; i < b.N; i++ {
			ips.AddString(loggen.PartiallyRandomIPv4(2))
		}
		mystore.AddDocument(fmt.Sprintf("test-%d.txt", doc), *ips)
	}
}

func BenchmarkStoreLeveldb(b *testing.B) {
	runStoreBench(b, "leveldb")
}
