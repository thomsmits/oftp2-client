GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

OBJECTS = bin/oftp2

SOURCE_FILES = cmd/oftp2/main.go \
			$(wildcard cmd/oftp2/*/*.go) \
			$(wildcard internal/liboftp2/*/*.go)

.PHONY: clean test test-coverage all

all: $(OBJECTS)

bin/oftp2: $(SOURCE_FILES)
	umask 022
	@mkdir -p $(@D)
	go build -o bin/ ./...

test: $(TEST_OBJECTS)
	go test  ./...

test-coverage: $(TEST_OBJECTS)
	go test -coverprofile cover.out ./...
	go tool cover -func=cover.out | grep "total:"

clean:
	rm -rf bin/
