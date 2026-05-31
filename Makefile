BINARY := yrank

.PHONY: build test vet clean

build:
	go build -o $(BINARY) .

test:
	go test -race ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY) coverage.txt
