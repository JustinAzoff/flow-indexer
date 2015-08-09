package ipset

import (
	"fmt"
	"net"
)

type Set struct {
	Store map[string]struct{}
}

func NewIpset() *Set {
	store := make(map[string]struct{})
	return &Set{store}
}

func (set *Set) AddString(s string) error {
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("Invalid IP Address %s", s)
	}

	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	keyString := string([]byte(ip))

	set.Store[keyString] = struct{}{}

	return nil
}
