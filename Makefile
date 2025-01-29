default: build

IMAGE := owenrumney/squealer

.PHONY: build
build: test
	./scripts/build.sh

.PHONY: test
test:
	go test -v -covermode=atomic -coverpkg ./... -coverprofile coverage.txt ./...


.PHONY: image
image:
	docker build --build-arg squealer_version=development -t $(IMAGE) .

.PHONY: lint
quality:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	golangci-lint run --timeout 3m --verbose