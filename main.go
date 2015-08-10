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
	"time"
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
	start := time.Now()
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
		}
	}
	duration := time.Since(start)
	fmt.Printf("%s: Read %d lines in %s\n", filename, lines, duration)

	start = time.Now()
	store.AddDocument(filename, *s)
	duration = time.Since(start)
	fmt.Printf("%s: Wrote %d unique ips in %s\n", filename, len(s.Store), duration)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	dbfile := os.Args[1]
	bs, err := NewLevelDBStore(dbfile)
	check(err)
	defer bs.Close()
	isFile := true
	arg := os.Args[2]
	if _, err := os.Stat(arg); os.IsNotExist(err) {
		isFile = false
	}
	if isFile {
		indexBroLog(bs, arg)
	} else {
		bs.QueryString(arg)
	}
}
