#!/bin/bash
BINARY=squealer
CHECK_GEN_BINARY=squealer-checkgen
TAG=${TRAVIS_TAG:-development}
GO111MODULE=on
export CGO_ENABLED=0
export GOFLAGS=-mod=vendor
args=(-ldflags "-X github.com/owenrumney/squealer/version.Version=${TAG} -s -w -extldflags '-fno-PIC -static'")

GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY}-darwin-amd64 "${args[@]}" ./cmd/squealer/
GOOS=linux GOARCH=amd64 go build -o bin/${BINARY}-linux-amd64 "${args[@]}" ./cmd/squealer/
GOOS=windows GOARCH=amd64 go build -o bin/${BINARY}-windows-amd64.exe "${args[@]}" ./cmd/squealer/