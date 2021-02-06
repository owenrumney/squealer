default: build

GOFMT_FILES?     = $$(find ./ -name '*.go' | grep -v vendor)
GOLANGCI_VERSION = v1.30.0
PKGER_VERSION    = 0.17.1
PLATFORM         := $(shell uname)

.PHONY: build
build: test
	./scripts/build.sh

.PHONY: test
test:
	GOSUMDB=off go test -v -mod=vendor ./...
