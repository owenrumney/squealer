default: build

IMAGE := owenrumney/squealer

.PHONY: build
build: test
	./scripts/build.sh

.PHONY: test
test:
	go test -v -covermode=atomic -coverpkg ./... -coverprofile coverage.txt ./...

.PHONY: push-image
push-image:
	./scripts/publish-image.sh

.PHONY: image
image:
	docker build --build-arg squealer_version=$(TRAVIS_TAG) -t $(IMAGE) .

.PHONY: quality
quality:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
	golangci-lint run --timeout 3m --verbose