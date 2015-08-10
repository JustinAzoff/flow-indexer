package main

import (
	"github.com/justinazoff/flow-indexer/ipset"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	bs, err := NewBoltStore("my.db")
	check(err)
	defer bs.Close()

	s := ipset.New()
	s.AddString("1.2.3.4")
	s.AddString("5.6.7.8")

	bs.AddDocument("flows-1.txt", *s)

	s2 := ipset.New()
	s2.AddString("1.2.3.4")
	s2.AddString("9.10.11.12")
	s2.AddString("2600::1")

	bs.AddDocument("flows-2.txt", *s2)

	bs.ListDocuments()
	bs.QueryString("1.2.3.4")

}
