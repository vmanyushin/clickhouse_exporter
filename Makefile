# Go parameters
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOGET=go get
BINARY_NAME=clickhouse_exporter
BINARY_UNIX=$(BINARY_NAME)_unix
VERSION ?= 1.0.0

all: dep test build
dep:
ifeq (, $(shell which dep))
	$(shell go get -u github.com/golang/dep/cmd/dep)
endif
	dep ensure
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v
test:
	$(GOTEST) $(go list ./... | grep -v /vendor/)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./...
	bin/$(BINARY_NAME)

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v
docker:
	docker build -t vmanyushin/${BINARY_NAME}:$(VERSION) .
	docker image push vmanyushin/${BINARY_NAME}:$(VERSION)
