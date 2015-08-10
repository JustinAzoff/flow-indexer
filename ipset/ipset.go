package ipset

import (
	"fmt"
	"net"
)

type Set struct {
	Store map[string]struct{}
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

func (set *Set) AddString(s string) error {
	keyString, err := IPToByteString(s)
	if err != nil {
		return err
	}
	set.Store[keyString] = struct{}{}
	return nil
}
