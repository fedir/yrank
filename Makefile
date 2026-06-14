BINARY := yrank
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build test vet clean coverage install release snapshot local-filter

# Thresholds for `local-filter` (0 = no limit on that dimension).
MIN_VIEWS ?= 0
MIN_LENGTH ?= 0

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

# Filter an existing CSV export by min views / min duration (seconds), locally, with
# no API quota. Reads IN, writes a filtered copy to OUT in the same CSV format:
#   make local-filter IN=sample_output/foo.csv OUT=foo_filtered.csv MIN_VIEWS=1000 MIN_LENGTH=300
local-filter: build
	@test -n "$(IN)" || { echo "IN is required, e.g. make local-filter IN=sample_output/foo.csv OUT=out.csv MIN_LENGTH=300"; exit 1; }
	@test -n "$(OUT)" || { echo "OUT is required, e.g. make local-filter IN=sample_output/foo.csv OUT=out.csv MIN_LENGTH=300"; exit 1; }
	./$(BINARY) -in $(IN) -out $(OUT) -min-views $(MIN_VIEWS) -min-length $(MIN_LENGTH)

clean:
	rm -f $(BINARY) coverage.txt
	rm -rf dist
