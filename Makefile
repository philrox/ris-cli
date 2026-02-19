BINARY := risgo

.PHONY: build test lint fmt clean install check

build:
	go build -o $(BINARY) .

test:
	go test -v -race ./...

lint:
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY)

install:
	go install .

check: fmt lint test
