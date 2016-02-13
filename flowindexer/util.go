package flowindexer

import (
	"fmt"
	"regexp"
)

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
