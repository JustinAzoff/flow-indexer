package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/JustinAzoff/flow-indexer/store"
	"github.com/spf13/cobra"
)

var bind string

var cmdWeb = &cobra.Command{
	Use:   "web [args]",
	Short: "Start http API",
	Long:  "Start http API",
	Run: func(cmd *cobra.Command, args []string) {
		startWeb()
	},
}

func init() {
	cmdWeb.Flags().StringVarP(&bind, "bind", "b", "127.0.0.1:8080", "Address to bind to")
	RootCmd.AddCommand(cmdWeb)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type storeHandler struct {
	stores []store.IpStore
}

func (sh *storeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, "Missing parameter: q", http.StatusBadRequest)
		return
	}
	for _, store := range sh.stores {
		docs, err := store.QueryString(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, doc := range docs {
			_, err := fmt.Fprintln(w, doc)
			if err != nil {
				return
			}
		}
	}
}

func startWeb() {
	var stores []store.IpStore
	matches, err := filepath.Glob("*.db")
	check(err)
	for _, db := range matches {
		log.Printf("Opening %s\n", db)
		mystore, err := store.NewStore("leveldb", db)
		check(err)
		stores = append(stores, mystore)
	}
	http.Handle("/search", &storeHandler{stores: stores})
	log.Fatal(http.ListenAndServe(bind, nil))
}
