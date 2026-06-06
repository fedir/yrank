BINARY := yrank
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build test vet clean coverage install release snapshot

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test -race ./...

vet:
	go vet ./...

coverage:
	go test -race -coverprofile=coverage.txt ./... && go tool cover -func=coverage.txt

install:
	go install $(LDFLAGS) .

# Build and publish a release for the current tag (used by CI on tag push).
release:
	goreleaser release --clean

# Build a local snapshot (no publish) to verify the release config.
snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY) coverage.txt
	rm -rf dist
