package ipset

import (
	"fmt"
	"net"
	"sort"
)

type Set struct {
	store map[string]struct{}
}

func New() *Set {
	store := make(map[string]struct{})
	return &Set{store}
}

func IPToByteString(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("Invalid IP Address %s", s)
	}

	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	return string([]byte(ip)), nil
}

func CIDRToByteStrings(s string) (string, string, error) {
	_, nw, err := net.ParseCIDR(s)
	if err != nil {
		return "", "", err
	}
	firstIP := nw.IP
	lastIP := make(net.IP, len(nw.IP))
	for i := 0; i < len(lastIP); i++ {
		lastIP[i] = firstIP[i] | ^nw.Mask[i]
	}
	return string([]byte(firstIP)), string([]byte(lastIP)), nil
}

func (set *Set) AddString(s string) error {
	keyString, err := IPToByteString(s)
	if err != nil {
		return err
	}
	set.store[keyString] = struct{}{}
	return nil
}

func (set *Set) Count() int {
	return len(set.store)
}

func (set *Set) SortedStrings() []string {
	strings := make([]string, set.Count())
	var i int = 0
	for ip, _ := range set.store {
		strings[i] = ip
		i++
	}
	sort.Strings(strings)
	return strings
}
