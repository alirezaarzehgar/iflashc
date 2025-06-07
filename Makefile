SRC := $(shell find -name "*.go")
BIN := iflashc
VERSION := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

all: build

build: ${BIN}

${BIN}: ${SRC}
		go build $(LDFLAGS) -o $@ .

install:
	install -oroot -groot -m 0775 ${BIN} /usr/bin

sqlc-gen:
	sqlc generate


clean:
	rm -f $(BIN)

.PHONY: build install clean sqlc-gen all
