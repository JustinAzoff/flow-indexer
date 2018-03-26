//Package ipset implements a set type for keeping track of unique ip addreses
package ipset

import (
	"fmt"
	"net"
	"sort"
)

//Set stores the unique ip addresses
//strings are used throughout because one cannot create a map[[]byte]..
type Set struct {
	store map[string]struct{}
}

//New returns an empty Set
func New() *Set {
	store := make(map[string]struct{})
	return &Set{store}
}

//IPToByteString converts an ip address to a byte array
func IPToByteString(ip net.IP) (string, error) {
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	return string([]byte(ip)), nil
}

//IPToByteString converts an ip address in string form to a byte array
func IPStringToByteString(s string) (string, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return "", fmt.Errorf("Invalid IP Address %q", s)
	}
	return IPToByteString(ip)
}

//CIDRToByteStrings converts a cidr block to its starting and ending addresses
//CIDRToByteStrings("192.168.1.0/24") will return
//IPToByteString("192.168.1.0"), IPToByteString("192.168.1.255"), error
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

//AddString adds a single ip address in string form into the Set
func (set *Set) AddString(s string) error {
	keyString, err := IPStringToByteString(s)
	if err != nil {
		return err
	}
	set.store[keyString] = struct{}{}
	return nil
}

//AddIP adds a single ip address in net.IP form into the Set
func (set *Set) AddIP(ip net.IP) error {
	keyString, err := IPToByteString(ip)
	if err != nil {
		return err
	}
	_, exists := set.store[keyString]
	if !exists {
		set.store[keyString] = struct{}{}
	}
	return nil
}

//Count returns the number of unique ip addresses in the Set
func (set *Set) Count() int {
	return len(set.store)
}

//SortedStrings returns all of the ip addresses in the Set in sorted order
func (set *Set) SortedStrings() []string {
	strings := make([]string, set.Count())
	var i int
	for ip := range set.store {
		strings[i] = ip
		i++
	}
	sort.Strings(strings)
	return strings
}

//SortedIPs returns all of the ip addresses in the Set in sorted order as IP objects
func (set *Set) SortedIPs() []net.IP {
	strings := set.SortedStrings()
	ips := make([]net.IP, len(strings))
	var i int
	for _, ip := range strings {
		ips[i] = net.IP(ip)
		i++
	}
	return ips
}
