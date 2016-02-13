package flowindexer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JustinAzoff/flow-indexer/store"
)

type indexerConfig struct {
	Name                    string `json:"name"`
	Backend                 string `json:"backend"`
	Indexer                 string `json:"indexer"`
	FileGlob                string `json:"file_glob"`
	FilenameToDatabaseRegex string `json:"filename_to_database_regex"`
	DatabaseRoot            string `json:"database_root"`
	DatabasePath            string `json:"datbase_path"`
}

type Config struct {
	Indexers []indexerConfig `json:"indexers"`
}

type Indexer struct {
	config indexerConfig
	stores []store.IpStore
}

type FlowIndexer struct {
	indexers map[string]Indexer
	config   Config
}

func loadConfig(filename string) (Config, error) {
	var cfg Config
	file, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg)
	return cfg, err
}

func NewFlowIndexerFromConfigFilename(filename string) (*FlowIndexer, error) {
	cfg, err := loadConfig(filename)
	if err != nil {
		return nil, err
	}
	indexerMap := make(map[string]Indexer)
	for _, indexercfg := range cfg.Indexers {
		indexer := Indexer{config: indexercfg}
		indexerMap[indexercfg.Name] = indexer
	}
	return &FlowIndexer{config: cfg, indexers: indexerMap}, nil
}

func (fi *FlowIndexer) GetIndexer(name string) (*Indexer, error) {
	indexer, ok := fi.indexers[name]
	if !ok {
		return nil, fmt.Errorf("Indexer %q not found", name)
	}
	return &indexer, nil
}

func (i *Indexer) ListDatabases() ([]string, error) {
	globPath := filepath.Join(i.config.DatabaseRoot, "*.db")
	databases, err := filepath.Glob(globPath)
	if err != nil {
		return databases, err
	}
	return databases, nil
}

func (i *Indexer) InitStores() error {
	databases, err := i.ListDatabases()
	if err != nil {
		return err
	}
	var stores []store.IpStore
	for _, db := range databases {
		log.Printf("Opening %s\n", db)
		mystore, err := store.NewStore("leveldb", db)
		if err != nil {
			return err
		}
		stores = append(stores, mystore)
	}
	i.stores = stores
	return nil
}

func RunDaemon(config string) {
	fi, err := NewFlowIndexerFromConfigFilename(config)
	if err != nil {
		log.Fatal(err)
	}

	indexer, err := fi.GetIndexer("bro")
	if err != nil {
		log.Fatal(err)
	}
	indexer.InitStores()
	for _, store := range indexer.stores {
		docs, err := store.QueryString("76.20.248.132")
		if err != nil {
			log.Fatal(err)
		}
		for _, doc := range docs {
			fmt.Printf("%s\n", doc)
		}
	}

}
