all: build test
build:
	go get -tags=embed github.com/tecbot/gorocksdb
	go get -t -v ./...
	go build -tags=embed
test:
	go test -tags=embed -v ./...

.PHONY: rpm
rpm: build
rpm: VERSION=$(shell ./flow-indexer version)
rpm:
	fpm -f -s dir -t rpm -n flow-indexer -v $(VERSION) \
	--iteration=1 \
	--architecture native \
	--description "Flow Indexer" \
	./flow-indexer=/usr/bin/
