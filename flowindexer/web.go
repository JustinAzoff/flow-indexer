package flowindexer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type fiHandler struct {
	fi *FlowIndexer
}

func (fh *fiHandler) handleIndexers(w http.ResponseWriter, req *http.Request) {
	var indexers []string
	for ix := range fh.fi.indexers {
		indexers = append(indexers, ix)
	}
	json.NewEncoder(w).Encode(indexers)
}
func (fh *fiHandler) handleSearch(w http.ResponseWriter, req *http.Request) {
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

	docs, err := indexer.QueryString(query)
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
func (fh *fiHandler) handleStats(w http.ResponseWriter, req *http.Request) {
	indexerParam := req.FormValue("i")
	query := req.FormValue("q")
	bucket := req.FormValue("bucket")
	if indexerParam == "" {
		http.Error(w, "Missing parameter: i", http.StatusBadRequest)
		return
	}
	if query == "" {
		http.Error(w, "Missing parameter: q", http.StatusBadRequest)
		return
	}
	bp, err := parseBucketParam(bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexer, err := fh.fi.GetIndexer(indexerParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := indexer.Stats(query, bp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (fh *fiHandler) handleExpandCIDR(w http.ResponseWriter, req *http.Request) {
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

	ips, err := indexer.ExpandCIDR(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, ip := range ips {
		_, err := fmt.Fprintln(w, ip)
		if err != nil {
			return
		}
	}
}
func (fh *fiHandler) handleDump(w http.ResponseWriter, req *http.Request) {
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

	err = indexer.Dump(query, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func startWeb(fi *FlowIndexer) {
	fh := &fiHandler{fi: fi}
	http.HandleFunc("/indexers", fh.handleIndexers)
	http.HandleFunc("/search", fh.handleSearch)
	http.HandleFunc("/stats", fh.handleStats)
	http.HandleFunc("/expandcidr", fh.handleExpandCIDR)
	http.HandleFunc("/dump", fh.handleDump)

	http.HandleFunc("/v1/indexers", fh.handleIndexers)
	http.HandleFunc("/v1/search", fh.handleSearch)
	http.HandleFunc("/v1/stats", fh.handleStats)
	http.HandleFunc("/v1/expandcidr", fh.handleExpandCIDR)
	http.HandleFunc("/v1/dump", fh.handleDump)

	bind := fi.config.HTTP.Bind
	log.Printf("Listening on %q\n", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
