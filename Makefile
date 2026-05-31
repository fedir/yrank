BINARY := yrank

.PHONY: build test vet clean coverage install

build:
	go build -o $(BINARY) .

test:
	go test -race ./...

vet:
	go vet ./...

coverage:
	go test -race -coverprofile=coverage.txt ./... && go tool cover -func=coverage.txt

install:
	go install .

clean:
	rm -f $(BINARY) coverage.txt
