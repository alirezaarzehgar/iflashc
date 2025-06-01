SRC := $(shell find -name "*.go")
BIN := iflashc

build:
	go build -o ${BIN} .

install:
	install -oroot -groot -m 0775 ${BIN} /usr/bin

.PHONY: build install
