    flow-indexer indexes flows

    Usage: 
      flow-indexer [command]

    Available Commands: 
      compact     Compact the database
      daemon      Start daemon
      expandcidr  Expand a CIDR range from those seen in the database
      index       Index flows
      search      Search flows
      help        Help about any command

    Flags:
          --dbpath="flows.db": Database path
      -h, --help[=false]: help for flow-indexer


    Use "flow-indexer [command] --help" for more information about a command.

Quickstart
==========

Create configuration
--------------------

    $ cp example_config.json config.json
    $ vi config.json

The indexer configuration is as follows:

* name - The name of the indexer. Keep this short and lowercase, as you will use it as an http query param.
* backend - The backend log ip extractor to use. Choices: bro, bro\_json, syslog.
* file\_glob - The shell globbing pattern that should match all of your log files.
* filename\_to\_database\_regex - A regular expression applied to each filename used to extract information used to name the database.
* database\_root - Where databases will be written to.  Should be indexer specific.
* datapath\_path - The name of an individual database.  This can contain $variables set in filename\_to\_database\_regex.

Adjust log paths and database paths.  The deciding factor for how to partition the databases is how many unique ips you see per day.
I suggest starting with monthly indexes.  If the performance takes a huge hit by the end of the month, switch to daily indexes.

    $ ./flow-indexer daemon

Query API
---------

    $ curl -s 'localhost:8080/search?i=conn&q=1.2.3.0/24"
    $ curl -s 'localhost:8080/stats?i=conn&q=1.2.3.0/24"

Lower level commands example
============================

Index flows
-----------

    ./flow-indexer --dbpath /tmp/f/flows.db index /tmp/f/conn*
    2016/02/06 23:36:51 /tmp/f/conn.00:00:00-01:00:00.log.gz: Read 4260 lines in 24.392765ms
    2016/02/06 23:36:51 /tmp/f/conn.00:00:00-01:00:00.log.gz: Wrote 281 unique ips in 2.215219ms
    2016/02/06 23:36:51 /tmp/f/conn.01:00:00-02:00:00.log.gz: Read 4376 lines in 24.186168ms
    2016/02/06 23:36:51 /tmp/f/conn.01:00:00-02:00:00.log.gz: Wrote 310 unique ips in 1.495277ms
    [...]
    2016/02/06 23:36:51 /tmp/f/conn.22:00:00-23:00:00.log.gz: Read 7799 lines in 18.350788ms
    2016/02/06 23:36:51 /tmp/f/conn.22:00:00-23:00:00.log.gz: Wrote 775 unique ips in 5.155262ms
    2016/02/06 23:36:51 /tmp/f/conn.23:00:00-00:00:00.log.gz: Read 5255 lines in 15.296847ms
    2016/02/06 23:36:51 /tmp/f/conn.23:00:00-00:00:00.log.gz: Wrote 400 unique ips in 2.910344ms

Re-Index flows
--------------

    ./flow-indexer --dbpath /tmp/f/flows.db index /tmp/f/conn*
    2016/02/06 23:37:36 /tmp/f/conn.00:00:00-01:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.01:00:00-02:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.02:00:00-03:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.03:00:00-04:00:00.log.gz Already indexed
    [...]
    2016/02/06 23:37:36 /tmp/f/conn.20:00:00-21:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.21:00:00-22:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.22:00:00-23:00:00.log.gz Already indexed
    2016/02/06 23:37:36 /tmp/f/conn.23:00:00-00:00:00.log.gz Already indexed

Expand CIDR Range
-----------------

    ./flow-indexer --dbpath /tmp/f/flows.db expandcidr 192.30.252.0/24
    192.30.252.86
    192.30.252.87
    192.30.252.92
    192.30.252.124
    192.30.252.125
    192.30.252.126
    192.30.252.127
    192.30.252.128
    192.30.252.129
    192.30.252.130
    192.30.252.131
    192.30.252.141

Search
------

    ./flow-indexer --dbpath /tmp/f/flows.db search 192.30.252.0/24
    /tmp/f/conn.03:00:00-04:00:00.log.gz
    /tmp/f/conn.04:00:00-05:00:00.log.gz
    /tmp/f/conn.06:00:00-07:00:00.log.gz
    /tmp/f/conn.14:00:00-15:00:00.log.gz
    /tmp/f/conn.18:00:00-19:00:00.log.gz
    /tmp/f/conn.20:00:00-21:00:00.log.gz
    /tmp/f/conn.22:00:00-23:00:00.log.gz

