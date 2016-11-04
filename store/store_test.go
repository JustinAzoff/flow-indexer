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
	{"1.2.3.3/32", []string{"/log/1.txt", "/log/2.txt"}},
	{"2.0.0.0/8", []string{"/log/2.txt"}},
	{"102:304::1", []string{"/log/3.txt"}},
}

var specialSearchTable = []struct {
	query string
	docs  []string
}{
	{"100.111.99.0/24", []string{"/log/special.txt"}},
	{"646f:633a::/64", []string{"/log/special.txt"}},
}

var basicExpandCidrTable = []struct {
	query string
	ips   []net.IP
}{
	{"1.2.3.0/24", makeIps([]string{"1.2.3.1", "1.2.3.2", "1.2.3.3", "1.2.3.4", "1.2.3.255"})},
	{"1.0.0.0/8", makeIps([]string{"1.2.3.1", "1.2.3.2", "1.2.3.3", "1.2.3.4", "1.2.3.255"})},
	{"2.0.0.0/8", makeIps([]string{"2.0.0.2", "2.0.0.3"})},

	//Once convered, this is starts with \x01\x02\x03\x04 in hex
	{"102:304::/8", makeIps([]string{"102:304::1"})},

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
	ips.AddString("1.2.3.255")
	s.AddDocument("/log/1.txt", *ips)

	ips.AddString("2.0.0.2")
	ips.AddString("2.0.0.3")
	s.AddDocument("/log/2.txt", *ips)

	ips = ipset.New()
	ips.AddString("102:304::1")
	s.AddDocument("/log/3.txt", *ips)

	for _, tt := range basicSearchTable {
		matches, err := s.QueryString(tt.query)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(matches, tt.docs) {
			t.Errorf("store.QueryString(%s) => %#v, want %#v", tt.query, matches, tt.docs)
		}
	}

	for _, tt := range basicExpandCidrTable {
		matches, err := s.ExpandCIDR(tt.query)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%v", matches) != fmt.Sprintf("%v", tt.ips) {
			t.Errorf("store.ExpandCIDR(%s) => %v, want %v", tt.query, matches, tt.ips)
		}
	}

	ips = ipset.New()
	ips.AddString("100.111.99.58") //doc: in hex
	ips.AddString("646f:633a::1")  //doc: in hex
	s.AddDocument("/log/special.txt", *ips)
	for _, tt := range specialSearchTable {
		matches, err := s.QueryString(tt.query)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(matches, tt.docs) {
			t.Errorf("store.QueryString(%s) => %#v, want %#v", tt.query, matches, tt.docs)
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

func runStoreBench(b *testing.B, storeType string, documents int) {
	mystore, err := NewStore(storeType, "test.db")
	if err != nil {
		b.Error(err)
		return
	}
	defer os.RemoveAll("test.db")
	for n := 0; n < b.N; n++ {
		for doc := 0; doc < documents; doc++ {
			b.StopTimer()
			ips := ipset.New()
			for i := 0; i < 30000; i++ {
				ips.AddString(loggen.PartiallyRandomIPv4(2))
			}
			b.StartTimer()
			mystore.AddDocument(fmt.Sprintf("test-%d-%d.txt", doc, n), *ips)
		}
	}
}

func BenchmarkStoreLeveldbDocs1(b *testing.B) {
	runStoreBench(b, "leveldb", 1)
}
func BenchmarkStoreLeveldbDocs24(b *testing.B) {
	runStoreBench(b, "leveldb", 24)
}
func BenchmarkStoreLeveldbDocs720(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping large store test in short mode")
	}

	runStoreBench(b, "leveldb", 720)
}
