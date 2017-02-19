all: build test
build:
	go get -t -v ./...
	go build
test:
	go test -tags=embed -v ./...
build_rocks:
	go get -tags=embed github.com/tecbot/gorocksdb
	go get -t -v ./...
	go build -tags='embed rocksdb'
test_rocks:
	go test -tags='embed rocksdb' -v ./...

.PHONY: rpm
rpm: build
rpm: VERSION=$(shell ./flow-indexer version)
rpm:
	fpm -f -s dir -t rpm -n flow-indexer -v $(VERSION) \
	--iteration=1 \
	--architecture native \
	--description "Flow Indexer" \
	./flow-indexer=/usr/bin/
