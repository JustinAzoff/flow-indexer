all: build test
build:
	go build
test:
	go test -v ./...

.PHONY: rpm
rpm: build
rpm: VERSION=$(shell ./flow-indexer version)
rpm:
	fpm -f -s dir -t rpm -n flow-indexer -v $(VERSION) \
	--iteration=1 \
	--architecture native \
	--description "Flow Indexer" \
	./flow-indexer=/usr/bin/
