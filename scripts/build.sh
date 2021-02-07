#!/bin/bash
BINARY=squealer
CHECK_GEN_BINARY=squealer-checkgen
TAG=${TRAVIS_TAG:-development}
GO111MODULE=on
export CGO_ENABLED=0
export GOFLAGS=-mod=vendor
args=(-ldflags "-X github.com/owenrumney/squealer/version.Version=${TAG} -s -w -extldflags '-fno-PIC -static'")

mkdir -p bin/darwin
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/${BINARY}-darwin-amd64 "${args[@]}" ./cmd/squealer/
mkdir -p bin/linux
GOOS=linux GOARCH=amd64 go build -o bin/linux/${BINARY}-linux-amd64 "${args[@]}" ./cmd/squealer/
mkdir -p bin/windows
GOOS=windows GOARCH=amd64 go build -o bin/windows/${BINARY}-windows-amd64.exe "${args[@]}" ./cmd/squealer/