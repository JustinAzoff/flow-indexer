package backend

import (
	"testing"
)

var expectedNFDUMPCount = 1241
var nfdumpTestFile = "test_data/nfdump.data"

func TestNFDUMPExtractIps(t *testing.T) {
	tests := []string{
		"nfdump-csv",
		"nfdump",
	}
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			ips, err := ExtractIps(tc, nfdumpTestFile)
			if err != nil {
				t.Fatal(err)
			}
			if ips.Count() != expectedNFDUMPCount {
				t.Errorf("ExtractIps('%s', '%s') count => %#v, want %#v", tc, nfdumpTestFile,
					ips.Count(), expectedNFDUMPCount)
			}

		})
	}
}

func TestNFDUMPImplementationsMatch(t *testing.T) {
	csv, err := ExtractIps("nfdump-csv", nfdumpTestFile)
	if err != nil {
		t.Fatal(err)
	}
	formatted, err := ExtractIps("nfdump", nfdumpTestFile)
	if err != nil {
		t.Fatal(err)
	}

	csvIPs := csv.SortedIPs()
	formattedIPs := formatted.SortedIPs()

	if len(csvIPs) != len(formattedIPs) {
		t.Fatalf("nfdump-csv count => %d, nfdump count => %d", len(csvIPs), len(formattedIPs))
	}
	a := csvIPs
	b := formattedIPs
	for idx := range csvIPs {
		if !a[idx].Equal(b[idx]) {
			t.Errorf("csvIPs[%d] != formattedIPs[%d] => %s != %s", idx, idx, a[idx], b[idx])
		}
	}

}

func BenchmarkNFDUMPCSV(b *testing.B) {
	benchmarks := []string{
		"nfdump-csv",
		"nfdump",
	}
	for _, bm := range benchmarks {
		b.Run(bm, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ExtractIps(bm, "test_data/nfdump.data")
			}
		})
	}
}
