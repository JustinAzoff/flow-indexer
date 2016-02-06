package main

import (
	"log"
	"os"
	"runtime"

	"github.com/justinazoff/flow-indexer/backend"
	"github.com/justinazoff/flow-indexer/store"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	dbfile := os.Args[1]
	bs, err := store.NewStore("leveldb", dbfile)
	check(err)
	defer bs.Close()
	isFile := true
	arg := os.Args[2]
	if _, err := os.Stat(arg); os.IsNotExist(err) {
		isFile = false
	}
	if isFile {
		backend.IndexBroLog(bs, arg)
	} else {
		bs.QueryString(arg)
	}
}
