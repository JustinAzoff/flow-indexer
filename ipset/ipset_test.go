package ipset

import (
	"testing"
)

var basicAddStringTests = []struct {
	in  string
	out string
}{
	{"1.2.3.4", "\x01\x02\x03\x04"},
	{"1.0.3.4", "\x01\x00\x03\x04"},
	{"2600::1", "\x26\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01"},
}

func TestIpsetAdd(t *testing.T) {
	for _, tt := range basicAddStringTests {
		s := NewIpset()
		s.AddString(tt.in)
		for k, _ := range s.Store {
			if k != tt.out {
				t.Errorf("Ipset.AddString(%#v) => %#v, want %#v", tt.in, k, tt.out)
			}
		}
	}
}
