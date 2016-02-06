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
	set.Store[keyString] = struct{}{}
	return nil
}
