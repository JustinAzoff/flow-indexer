package flowindexer

import (
	"fmt"
	"github.com/JustinAzoff/flow-indexer/store"
)

func RunSearch(dbpath string, args []string) {
	mystore, err := store.NewStore("leveldb", dbpath)
	check(err)
	defer mystore.Close()

	docs, err := mystore.QueryString(args[0])
	check(err)
	for _, doc := range docs {
		fmt.Println(doc)
	}
}
