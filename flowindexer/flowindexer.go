package flowindexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/JustinAzoff/flow-indexer/backend"
	"github.com/JustinAzoff/flow-indexer/ipset"
	"github.com/JustinAzoff/flow-indexer/store"
)

type indexerConfig struct {
	Name                    string `json:"name"`
	Backend                 string `json:"backend"`
	Indexer                 string `json:"indexer"`
	FileGlob                string `json:"file_glob"`
	Store                   string `json:"store"`
	FilenameToDatabaseRegex string `json:"filename_to_database_regex"`
	FilenameToTimeRegexp    string `json:"filename_to_time_regex"`
	DatabaseRoot            string `json:"database_root"`
	DatabasePath            string `json:"database_path"`
}

type httpConfig struct {
	Bind string `json:"bind"`
}

type Config struct {
	Indexers []indexerConfig `json:"indexers"`
	HTTP     httpConfig      `json:"http"`
}

type Indexer struct {
	config               indexerConfig
	stores               []store.IpStore
	storeMap             map[string]store.IpStore
	indexedFilenames     map[string]bool
	filenameToTimeRegexp *regexp.Regexp
}

type FlowIndexer struct {
	indexers map[string]*Indexer
	config   Config
}

type BucketHit struct {
	Bucket string `json:"bucket"`
	Hits   int    `json:"hits"`
}

type queryStat struct {
	Hits  int    `json:"hits"`
	First string `json:"first"`
	Last  string `json:"last"`

	FirstTime time.Time `json:"first_time"`
	LastTime  time.Time `json:"last_time"`

	Buckets []*BucketHit `json:"buckets"`
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
	indexerMap := make(map[string]*Indexer)
	for _, indexercfg := range cfg.Indexers {
		if indexercfg.Store == "" {
			indexercfg.Store = store.DefaultStore
		}
		indexer := Indexer{config: indexercfg}
		indexer.storeMap = make(map[string]store.IpStore)
		indexer.indexedFilenames = make(map[string]bool)
		if indexercfg.FilenameToTimeRegexp != "" {
			indexer.filenameToTimeRegexp = regexp.MustCompile(indexercfg.FilenameToTimeRegexp)
		}
		indexerMap[indexercfg.Name] = &indexer
	}
	return &FlowIndexer{config: cfg, indexers: indexerMap}
}

func (fi *FlowIndexer) GetIndexer(name string) (*Indexer, error) {
	indexer, ok := fi.indexers[name]
	if !ok {
		return nil, fmt.Errorf("Indexer %q not found", name)
	}
	return indexer, nil
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

func (i *Indexer) FilenameToTime(filename string) (time.Time, error) {
	tm, err := logFilenameToTime(filename, i.filenameToTimeRegexp)
	if err != nil {
		return time.Now(), err
	}
	return tm, nil
}

func (i *Indexer) IndexOne(filename string, checkGrowing bool) error {
	_, alreadyIndexed := i.indexedFilenames[filename]
	if alreadyIndexed {
		return nil
	}
	if checkGrowing {
		stillGrowing, err := isFileGrowing(filename)
		if err != nil {
			log.Printf("Failed to check if %q is growing: %q", filename, err)
			return err
		}
		if stillGrowing {
			log.Printf("Skipping still growing file %q\n", filename)
			return nil
		}
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
	for idx, l := range logs {
		//Assume the last file in the list is the most recent one
		//and check to see if it is still growing before indexing it
		checkGrowing := idx == len(logs)-1
		i.IndexOne(l, checkGrowing)
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
func (i *Indexer) QueryString(query string) ([]string, error) {
	var documents []string
	for _, store := range i.stores {
		docs, err := store.QueryString(query)
		if err != nil {
			return documents, err
		}
		documents = append(documents, docs...)
	}
	return documents, nil
}

func (i *Indexer) ExpandCIDR(query string) ([]net.IP, error) {
	allips := ipset.New()

	for _, store := range i.stores {
		ips, err := store.ExpandCIDR(query)
		if err != nil {
			return []net.IP{}, err
		}
		for _, ip := range ips {
			allips.AddIP(ip)
		}
	}
	return allips.SortedIPs(), nil
}

func (i *Indexer) Stats(query string) (queryStat, error) {
	var stat = queryStat{}
	docs, err := i.QueryString(query)
	sort.Strings(docs)
	if err != nil {
		return stat, err
	}
	if len(docs) > 0 {
		stat.First = docs[0]
		stat.Last = docs[len(docs)-1]

		if t, err := i.FilenameToTime(stat.First); err == nil {
			stat.FirstTime = t
		}
		if t, err := i.FilenameToTime(stat.Last); err == nil {
			stat.LastTime = t
		}
	}
	stat.Hits = len(docs)

	var last string
	var bp *BucketHit //a pointer to a bucket hit
	for _, doc := range docs {
		if t, err := i.FilenameToTime(doc); err == nil {
			bucket := t.Truncate(time.Hour * 24).Format(time.RFC3339)
			if bucket != last {
				bh := BucketHit{Bucket: bucket, Hits: 1}
				stat.Buckets = append(stat.Buckets, &bh)
				last = bucket
				bp = &bh
			} else {
				bp.Hits++
			}
		}
	}

	return stat, nil
}

func (i *Indexer) Dump(query string, writer io.Writer) error {
	docs, err := i.QueryString(query)
	if err != nil {
		return err
	}

	for _, fn := range docs {
		err = backend.FilterIPs(i.config.Backend, fn, query, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunIndexAll(config string) {
	fi, err := NewFlowIndexerFromConfigFilename(config)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup

	for _, indexer := range fi.indexers {
		wg.Add(1)
		go func(indexer *Indexer) {
			defer wg.Done()
			indexer.IndexAll()
		}(indexer)
	}
	wg.Wait()
}

func RunDaemon(config string) {
	fi, err := NewFlowIndexerFromConfigFilename(config)
	if err != nil {
		log.Fatal(err)
	}

	//Before starting the API, make sure all the stores are open
	for _, indexer := range fi.indexers {
		indexer.RefreshStores()
	}
	go startWeb(fi)

	for _, indexer := range fi.indexers {
		go func(indexer *Indexer) {
			for {
				indexer.RefreshStores()
				indexer.IndexAll()
				time.Sleep(60 * time.Second)
			}
		}(indexer)
	}
	for {
		time.Sleep(5 * time.Second)
	}

}
