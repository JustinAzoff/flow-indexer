package flowindexer

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func logFilenameToTime(filename string, filenameToTimeRegex *regexp.Regexp) (time.Time, error) {
	if filenameToTimeRegex == nil {
		return time.Now(), fmt.Errorf("filename_to_time_regex is not set in indexer configuration")
	}
	n1 := filenameToTimeRegex.SubexpNames()
	r2 := filenameToTimeRegex.FindAllStringSubmatch(filename, -1)[0]

	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}

	getOrZero := func(key string) int {
		val, exists := md[key]
		if !exists {
			return 0
		}
		num, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return num
	}
	year := getOrZero("year")
	month := time.Month(getOrZero("month"))
	day := getOrZero("day")
	hour := getOrZero("hour")
	minute := getOrZero("minute")
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC), nil
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

const (
	minuteFmt = "2006-01-02T15:04:05"
	hourFmt   = "2006-01-02T15"
	dayFmt    = "2006-01-02"
	monthFmt  = "2006-01"
	yearFmt   = "2006"
)

func timeToBucket(tm time.Time, trunc string) (string, error) {
	var bucket string
	switch trunc {
	case "minute":
		bucket = tm.Format(minuteFmt)
	case "hour":
		bucket = tm.Format(hourFmt)
	case "day":
		bucket = tm.Format(dayFmt)
	case "month":
		bucket = tm.Format(monthFmt)
	case "year":
		bucket = tm.Format(yearFmt)
	default:
		return bucket, fmt.Errorf("Invalid truncation period: %s", trunc)
	}
	return bucket, nil
}

func parseBucketParam(bucket string) (bucketParam, error) {

	if bucket == "" {
		return bucketParam{"month", "day"}, nil
	}

	bp := bucketParam{}
	parts := strings.Split(bucket, "/")
	switch len(parts) {
	case 1:
		bp.groupby = parts[0]
		switch bp.groupby {
		case "year":
			bp.count = "month"
		case "month":
			bp.count = "day"
		case "day":
			bp.count = "hour"
		}
	case 2:
		bp.groupby = parts[0]
		bp.count = parts[1]
	default:
		return bp, fmt.Errorf("Invalid bucketing specification %q. Only one / is allowed", bucket)
	}
	return bp, nil
}
