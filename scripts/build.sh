#!/bin/bash

TAG=${TRAVIS_TAG:-development}
GO111MODULE=on
export CGO_ENABLED=0
export GOFLAGS=-mod=vendor
args=(-ldflags "-X github.com/owenrum/squealer/version.Version=${TAG} -s -w -extldflags '-fno-PIC -static'")

GOOS=darwin GOARCH=amd64 go build -o "bin/squealer.darwin.amd64.${TAG}" "${args[@]}" ./cmd/squealer/
GOOS=linux GOARCH=amd64 go build -o "bin/squealer.linux.amd64.${TAG}" "${args[@]}" ./cmd/squealer/
GOOS=windows GOARCH=amd64 go build -o "bin/squealer.windows.amd64.${TAG}.exe" "${args[@]}" ./cmd/squealer/
