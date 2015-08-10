package main

import (
	"bufio"
	"fmt"
	"github.com/justinazoff/flow-indexer/ipset"
	gzip "github.com/klauspost/pgzip"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func indexBroLog(store *LevelDBStore, filename string) {
	exists, err := store.HasDocument(filename)
	if err != nil {
		fmt.Print(err)
		return
	}
	if exists {
		fmt.Printf("%s Already indexed\n", filename)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Print(err)
	}
	reader, err := gzip.NewReader(f)
	if err != nil {
		fmt.Print(err)
	}
	br := bufio.NewReader(reader)

	s := ipset.New()
	lines := 0
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print(err)
			return
		}
		if line[0] != '#' {
			parts := strings.Split(line, "\t")
			s.AddString(parts[2])
			s.AddString(parts[4])
			lines += 1
			if lines%1000 == 0 {
				fmt.Printf("\rRead %d lines", lines)
			}
		}
	}
	fmt.Printf("\rRead %d lines\n", lines)

	store.AddDocument(filename, *s)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	bs, err := NewLevelDBStore("my.db")
	check(err)
	defer bs.Close()

	arg := os.Args[1]
	isFile := true
	if _, err := os.Stat(arg); os.IsNotExist(err) {
		isFile = false
	}
	if isFile {
		indexBroLog(bs, arg)
	}
	bs.ListDocuments()
	if !isFile {
		bs.QueryString(arg)
	}
}
