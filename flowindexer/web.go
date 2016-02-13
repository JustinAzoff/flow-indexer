package flowindexer

import (
	"fmt"
	"log"
	"net/http"
)

type fiHandler struct {
	fi *FlowIndexer
}

func (fh *fiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	indexerParam := req.FormValue("i")
	query := req.FormValue("q")
	if indexerParam == "" {
		http.Error(w, "Missing parameter: i", http.StatusBadRequest)
		return
	}
	if query == "" {
		http.Error(w, "Missing parameter: q", http.StatusBadRequest)
		return
	}

	indexer, err := fh.fi.GetIndexer(indexerParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, store := range indexer.stores {
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

func startWeb(fi *FlowIndexer) {
	http.Handle("/search", &fiHandler{fi: fi})
	bind := fi.config.HTTP.Bind
	log.Printf("Listening on %q\n", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
