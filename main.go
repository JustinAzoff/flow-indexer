package main

import (
	"log"
	"os"
	"path/filepath"
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
	arg := os.Args[2]

	err = mystore.QueryString(arg)
	if err == nil {
		return
	}

	myindexer := backend.NewBackend("bro")
	matches, err := filepath.Glob(arg)
	check(err)
	for _, fp := range matches {
		err = Index(mystore, myindexer, fp)
	}
	check(err)
}
