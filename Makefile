BINARY_NAME=md
VERSION ?= 1.0.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -s -w"

.PHONY: all build test test-race vet clean install doctor demo

all: build

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/md

test:
	go test -v ./...

test-race:
	go test -race -v ./...

vet:
	go vet ./...

clean:
	rm -rf bin/

install: build
	cp bin/$(BINARY_NAME) $(shell go env GOPATH)/bin/$(BINARY_NAME)

doctor: build
	./bin/$(BINARY_NAME) doctor

demo: build
	./bin/$(BINARY_NAME) testdata/demo.md
