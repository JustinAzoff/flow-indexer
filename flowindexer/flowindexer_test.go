package flowindexer

import (
	"testing"
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

var testConfig = []byte(`
{
    "http": {
        "bind": ":8080"
    },
    "indexers": [ {
        "name": "bro",
        "backend": "bro",
        "file_glob": "/home/justin/tmp/bro_logs/conn.*",
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
