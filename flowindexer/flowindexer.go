package flowindexer

import (
	"fmt"
	"github.com/JustinAzoff/flow-indexer/store"
	"regexp"
)

type indexerConfig struct {
	Backend           string `json:"backend"`
	Indexer           string `json:indexer"`
	FileGlob          string `json:file_glob"`
	FilenameToDbRegex string `json:filename_to_db_regex`
	RegexReplacement  string `json:regex_replacement`
}

type FlowIndexer struct {
	stores []store.IpStore
	config indexerConfig
}

func logFilenameToDatabase(filename, filenameToDbRegex, regexReplacement string) (string, error) {
	re, err := regexp.Compile(filenameToDbRegex)
	if err != nil {
		return "", err
	}
	submatches := re.FindSubmatchIndex([]byte(filename))
	if submatches == nil {
		return "", fmt.Errorf("%q did not match %q", filename, filenameToDbRegex)
	}
	var db []byte
	db = re.ExpandString(db[:], regexReplacement, filename, submatches)

	dbString := string(db)
	return dbString, nil
}
