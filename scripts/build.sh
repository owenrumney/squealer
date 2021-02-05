#!/usr/bin/env bash

set -ex

if [[ ${GOOS} = "" ]]; then
    GOOS=$(uname | tr "[:upper:]" "[:lower:]")
fi

output=bin/squealer.${GOOS}.amd64.$TRAVIS_TAG
mkdir -p bin
env GOSUMDB=off go build -mod=vendor -o $output ./cmd/squealer