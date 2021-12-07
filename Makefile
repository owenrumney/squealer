default: build

IMAGE := owenrum/squealer

.PHONY: build
build: test
	./scripts/build.sh

.PHONY: test
test:
	go test -v -covermode=atomic -coverpkg ./... -coverprofile coverage.txt ./...

.PHONY: push-image
push-image:
	./scripts/publish-image.sh

image:
	docker build --build-arg squealer_version=$(TRAVIS_TAG) -t $(IMAGE) .