all: clean test build

build: clean test
	go build

clean:
	go clean ./...

test: clean
	go test ./...

install: clean test
	go install
