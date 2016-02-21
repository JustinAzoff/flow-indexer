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
		s := New()
		s.AddString(tt.in)
		for _, ip := range s.SortedStrings() {
			if ip != tt.out {
				t.Errorf("Ipset.AddString(%#v) => %#v, want %#v", tt.in, ip, tt.out)
			}
		}
	}
}

var basicCIDRToByteStrings = []struct {
	in       string
	startout string
	endout   string
}{
	{"1.2.3.4/16",
		"\x01\x02\x00\x00",
		"\x01\x02\xff\xff"},
	{"1.2.3.4/24",
		"\x01\x02\x03\x00",
		"\x01\x02\x03\xff"},
	{"fe80::a21e:13ff:ab31:3d60/64",
		"\xfe\x80\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
		"\xfe\x80\x00\x00\x00\x00\x00\x00\xff\xff\xff\xff\xff\xff\xff\xff"},
}

func TestCIDRToByteStrings(t *testing.T) {
	for _, tt := range basicCIDRToByteStrings {
		s, e, err := CIDRToByteStrings(tt.in)
		if err != nil {
			t.Error(err)
		}
		if s != tt.startout {
			t.Errorf("Ipset.CIDRToByteStrings(%#v) => start is %#v, want %#v", tt.in, s, tt.startout)
		}
		if e != tt.endout {
			t.Errorf("Ipset.CIDRToByteStrings(%#v) => end is %#v, want %#v", tt.in, e, tt.endout)
		}
	}
}
