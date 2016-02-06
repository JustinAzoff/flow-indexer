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
	mystore, err := store.NewStore("leveldb", dbfile)
	check(err)
	defer mystore.Close()
	isFile := true
	arg := os.Args[2]
	if _, err := os.Stat(arg); os.IsNotExist(err) {
		isFile = false
	}
	if isFile {
		myindexer := backend.NewBackend("bro")
		check(err)
		Index(mystore, myindexer, arg)
	} else {
		mystore.QueryString(arg)
	}
}
