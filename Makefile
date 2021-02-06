default: build

.PHONY: build
build: test
    ./scripts/build.sh

.PHONY: test
test:
    go test -v -covermode=atomic -coverpkg ./... -coverprofile coverage.txt ./...
