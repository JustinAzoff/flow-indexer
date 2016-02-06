package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/JustinAzoff/flow-indexer/cmd"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
