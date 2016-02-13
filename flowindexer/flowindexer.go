package flowindexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JustinAzoff/flow-indexer/backend"
	"github.com/JustinAzoff/flow-indexer/store"
)

type indexerConfig struct {
	Name                    string `json:"name"`
	Backend                 string `json:"backend"`
	Indexer                 string `json:"indexer"`
	FileGlob                string `json:"file_glob"`
	Store                   string `json:"store"`
	FilenameToDatabaseRegex string `json:"filename_to_database_regex"`
	DatabaseRoot            string `json:"database_root"`
	DatabasePath            string `json:"database_path"`
}

type Config struct {
	Indexers []indexerConfig `json:"indexers"`
}

type Indexer struct {
	config           indexerConfig
	stores           []store.IpStore
	storeMap         map[string]store.IpStore
	indexedFilenames map[string]bool
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

func parseConfig(jsonBlob []byte) (Config, error) {
	var cfg Config
	err := json.Unmarshal(jsonBlob, &cfg)
	return cfg, err
}

func NewFlowIndexerFromConfigFilename(filename string) (*FlowIndexer, error) {
	cfg, err := loadConfig(filename)
	if err != nil {
		return nil, err
	}
	return NewFlowIndexerFromConfig(cfg), nil
}
func NewFlowIndexerFromConfigBytes(jsonBlob []byte) (*FlowIndexer, error) {
	cfg, err := parseConfig(jsonBlob)
	if err != nil {
		return nil, err
	}
	return NewFlowIndexerFromConfig(cfg), nil
}
func NewFlowIndexerFromConfig(cfg Config) *FlowIndexer {
	indexerMap := make(map[string]Indexer)
	for _, indexercfg := range cfg.Indexers {
		if indexercfg.Store == "" {
			indexercfg.Store = store.DefaultStore
		}
		indexer := Indexer{config: indexercfg}
		indexer.storeMap = make(map[string]store.IpStore)
		indexer.indexedFilenames = make(map[string]bool)
		indexerMap[indexercfg.Name] = indexer
	}
	return &FlowIndexer{config: cfg, indexers: indexerMap}
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

func (i *Indexer) ListLogs() ([]string, error) {
	logs, err := filepath.Glob(i.config.FileGlob)
	if err != nil {
		return logs, err
	}
	return logs, nil
}

func (i *Indexer) FilenameToDatabaseFilename(filename string) (string, error) {
	db, err := logFilenameToDatabase(filename, i.config.FilenameToDatabaseRegex, i.config.DatabasePath)
	if err != nil {
		return "", err
	}
	if db == "" {
		return "", errors.New("Empty string return from logFilenameToDatabase")
	}
	return filepath.Join(i.config.DatabaseRoot, db), nil
}

func (i *Indexer) IndexOne(filename string) error {
	_, alreadyIndexed := i.indexedFilenames[filename]
	if alreadyIndexed {
		return nil
	}

	dbPath, err := i.FilenameToDatabaseFilename(filename)
	if err != nil {
		log.Printf("Can't convert %q to database filename: %q", filename, err)
		return err
	}
	s, err := i.OpenOrCreateStore(dbPath)
	if err != nil {
		log.Printf("Error opening database %q %q", dbPath, err)
		return err
	}
	mybackend := backend.NewBackend(i.config.Backend)
	err = Index(*s, mybackend, filename)
	if err != nil {
		log.Printf("Error indexing %q %q", filename, err)
		return err
	}
	i.indexedFilenames[filename] = true
	return nil
}

func (i *Indexer) IndexAll() error {
	logs, err := i.ListLogs()
	if err != nil {
		return err
	}
	for _, l := range logs {
		i.IndexOne(l)
	}

	return nil
}
func (i *Indexer) OpenOrCreateStore(filename string) (*store.IpStore, error) {
	s, alreadyExists := i.storeMap[filename]
	if alreadyExists {
		return &s, nil
	}
	log.Printf("Opening %s\n", filename)
	mystore, err := store.NewStore(i.config.Store, filename)
	if err != nil {
		return nil, err
	}
	i.storeMap[filename] = mystore
	i.stores = append(i.stores, mystore)
	return &mystore, nil
}

func (i *Indexer) RefreshStores() error {
	databases, err := i.ListDatabases()
	if err != nil {
		return err
	}

	seenDatabases := make(map[string]bool)
	for _, store := range i.stores {
		seenDatabases[store.Filename()] = false
	}

	for _, db := range databases {
		seenDatabases[db] = true
		_, alreadyExists := i.storeMap[db]
		if alreadyExists {
			continue
		}
		log.Printf("Opening %s\n", db)
		_, err = i.OpenOrCreateStore(db)
		if err != nil {
			log.Printf("Error opening %s: %s\n", db, err)
		}
	}
	//Now, see if any of our databases no longer exist on disk
	newStores := i.stores[:0]
	for _, s := range i.stores {
		if seenDatabases[s.Filename()] {
			newStores = append(newStores, s)
		} else {
			s.Close()
			log.Printf("Closing %s\n", s.Filename())
			delete(i.storeMap, s.Filename())
		}
	}
	i.stores = newStores
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

	for {

		indexer.RefreshStores()
		indexer.IndexAll()
		for _, store := range indexer.stores {
			log.Printf("Looking at store %s\n", store.Filename())
			docs, err := store.QueryString("76.20.248.132")
			if err != nil {
				log.Fatal(err)
			}
			for _, doc := range docs {
				fmt.Printf("%s\n", doc)
			}
		}
		time.Sleep(2 * time.Second)
	}

}
