SRC := $(shell find -name "*.go")
BIN := iflashc
VERSION := $(shell git describe --tags 2>/dev/null)
LDFLAGS := -ldflags="-s -w -X=github.com/alirezaarzehgar/iflashc/cmd.Version=${VERSION}"

all: download build

build: ${BIN}

download:
	go mod download

${BIN}: ${SRC}
	go build $(LDFLAGS) -o $@ .

install: build
	install -oroot -groot -m 0775 ${BIN} /usr/bin

sqlc-gen:
	sqlc generate


clean:
	rm -f $(BIN)

.PHONY: build install clean sqlc-gen all
