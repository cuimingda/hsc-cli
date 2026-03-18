GO ?= go

.PHONY: install test

install:
	$(GO) install ./cmd/hsc

test:
	$(GO) test ./...
