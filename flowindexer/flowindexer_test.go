package flowindexer

import (
	"regexp"
	"testing"
	"time"
)

var basiclogFilenameToDatabase = []struct {
	filename    string
	regex       string
	replacement string
	database    string
}{
	{"/bro/logs/2015-01-01/conn.blah.doesntmatter.log.gz", "logs/(?P<db>\\d+)-\\d+-\\d+", "$db.db", "2015.db"},
	{"/bro/logs/2015-01-01/conn.blah.doesntmatter.log.gz", "logs/(?P<db>\\d+-\\d+)-\\d+", "$db.db", "2015-01.db"},
	{"/bro/logs/2015-01-01/conn.blah.doesntmatter.log.gz", "logs/(?P<db>\\d+-\\d+-\\d+)", "$db.db", "2015-01-01.db"},
}

func TestLogFilenameToDatabase(t *testing.T) {
	for _, tt := range basiclogFilenameToDatabase {
		db, err := logFilenameToDatabase(tt.filename, tt.regex, tt.replacement)
		if err != nil {
			t.Error(err)
		}
		if db != tt.database {
			t.Errorf("flowindexer.logFilenameToDatabase(%#v) => db is %#v, want %#v", tt.filename, db, tt.database)
		} else {
			t.Logf("flowindexer.logFilenameToDatabase(%#v) => %#v", tt.filename, db)
		}
	}
}

var basiclogFilenameToTime = []struct {
	filename string
	regex    string
	tm       time.Time
}{
	{"/bro/logs/2016-01-01/conn.10:15:00-11:00:00.log.gz",
		"logs/(?P<year>\\d\\d\\d\\d)-(?P<month>\\d\\d)-(?P<day>\\d\\d)/\\w+\\.(?P<hour>\\d\\d):(?P<minute>\\d\\d)",
		time.Date(2016, 1, 1, 10, 15, 0, 0, time.UTC)},
}

func TestLogFilenameToTime(t *testing.T) {
	for _, tt := range basiclogFilenameToTime {
		re, err := regexp.Compile(tt.regex)
		if err != nil {
			t.Error(err)
		}
		tm, err := logFilenameToTime(tt.filename, re)
		if err != nil {
			t.Error(err)
		}
		if !tt.tm.Equal(tm) {
			t.Errorf("flowindexer.logFilenameToTime(%#v) => tm is %#v, want %#v", tt.filename, tm, tt.tm)
		} else {
			t.Logf("flowindexer.logFilenameToTime(%#v) => %#v", tt.filename, tm)
		}
	}
}

var timeToBucketTests = []struct {
	tm    time.Time
	trunc string
	out   string
}{
	{
		time.Date(2016, 2, 3, 10, 15, 0, 0, time.UTC),
		"hour",
		"2016-02-03T10",
	},
	{
		time.Date(2016, 2, 3, 10, 15, 0, 0, time.UTC),
		"day",
		"2016-02-03",
	},
	{
		time.Date(2016, 2, 3, 10, 15, 0, 0, time.UTC),
		"month",
		"2016-02",
	},
	{
		time.Date(2016, 2, 3, 10, 15, 0, 0, time.UTC),
		"year",
		"2016",
	},
}

func TestTimeToBucket(t *testing.T) {
	for _, tt := range timeToBucketTests {
		out, err := timeToBucket(tt.tm, tt.trunc)
		if err != nil {
			t.Error(err)
		}
		if tt.out != out {
			t.Errorf("flowindexer.timeToBucket(%q, %#v) => out is %#v, want %#v", tt.tm, tt.trunc, out, tt.out)
		} else {
			t.Logf("flowindexer.timeToBucket(%q, %#v) => %#v", tt.tm, tt.trunc, tt.out)
		}
	}
}

var testConfig = []byte(`
{
    "http": {
        "bind": ":8080"
    },
    "indexers": [ {
        "name": "bro",
        "backend": "bro",
        "file_glob": "/home/justin/tmp/bro_logs/conn.*",
		"filename_to_time_regex": "logs/(?P<year>\\d\\d\\d\\d)-(?P<month>\\d\\d)-(?P<day>\\d\\d)/\\w+\\.(?P<hour>\\d\\d):(?P<minute>\\d\\d)",
        "database_root": "/home/justin/tmp/bro_logs",
        "database_path": "db.db"
        }
    ]
}
`)

func TestNewFlowIndexerFromConfigBytes(t *testing.T) {
	fi, err := NewFlowIndexerFromConfigBytes(testConfig)
	if err != nil {
		t.Error(err)
	}
	if fi.indexers["bro"].config.Name != "bro" {
		t.Errorf("Something wrong with config")
	}
}

var testDocuments = []string{
	"/bro/logs/2015-04-08/conn.10:00:00-11:00:00.log.gz",
	"/bro/logs/2015-04-08/conn.11:00:00-12:00:00.log.gz",
	"/bro/logs/2016-05-04/conn.17:00:00-18:00:00.log.gz",
	"/bro/logs/2016-05-04/conn.18:00:00-19:00:00.log.gz",
	"/bro/logs/2016-05-05/conn.03:00:00-04:00:00.log.gz",
	"/bro/logs/2016-05-05/conn.05:00:00-06:00:00.log.gz",
}

type bucketTest struct {
	index  int
	bucket string
	hits   int
}

func checkBuckets(t *testing.T, which string, stats queryStat, tests []bucketTest) {
	for _, tt := range tests {
		if stats.Buckets[tt.index].Bucket != tt.bucket {
			t.Errorf("%s stats.Bucket[%d].Bucket => %q, want %q", which, tt.index, stats.Buckets[tt.index].Bucket, tt.bucket)
		}
		if stats.Buckets[tt.index].Hits != tt.hits {
			t.Errorf("%s stats.Bucket[%d].Hits => %d, want %d", which, tt.index, stats.Buckets[tt.index].Hits, tt.hits)
		}
	}
}

var testDocumentsBucketedByMonthHour = []bucketTest{
	{0, "2015-04", 2},
	{1, "2016-05", 4},
}
var testDocumentsBucketedByMonthDay = []bucketTest{
	{0, "2015-04", 1},
	{1, "2016-05", 2},
}

func TestFilenamesToStats(t *testing.T) {
	fi, err := NewFlowIndexerFromConfigBytes(testConfig)
	if err != nil {
		t.Error(err)
	}
	i := fi.indexers["bro"]
	stats, err := i.FilenamesToStats(testDocuments, "month", "hour")
	if err != nil {
		t.Error(err)
		return
	}
	if stats.First != "/bro/logs/2015-04-08/conn.10:00:00-11:00:00.log.gz" {
		t.Errorf("stats.First should not be %q", stats.First)
	}
	if stats.Last != "/bro/logs/2016-05-05/conn.05:00:00-06:00:00.log.gz" {
		t.Errorf("stats.Last should not be %q", stats.Last)
	}

	checkBuckets(t, "month/hour", stats, testDocumentsBucketedByMonthHour)

	stats, err = i.FilenamesToStats(testDocuments, "month", "day")
	if err != nil {
		t.Error(err)
		return
	}
	checkBuckets(t, "month/day", stats, testDocumentsBucketedByMonthDay)
}
