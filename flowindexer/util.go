package flowindexer

import (
	"fmt"
	"os"
	"regexp"
	"time"
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

func isFileGrowing(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}
	oldBytes := stat.Size()

	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)

		stat, err = file.Stat()
		if err != nil {
			return false, err
		}
		newBytes := stat.Size()

		if newBytes != oldBytes {
			return true, nil
		}
	}
	return false, nil
}
