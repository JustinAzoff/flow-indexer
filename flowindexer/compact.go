package flowindexer

import (
	"github.com/JustinAzoff/flow-indexer/store"
)

func RunCompact(dbpath string) {
	mystore, err := store.NewStore("leveldb", dbpath)
	check(err)
	defer mystore.Close()
	err = mystore.Compact()
	check(err)
}
