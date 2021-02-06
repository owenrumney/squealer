#!/bin/bash

set -eo pipefail

if [ -z "$GITHUB_TOKEN" ]; then echo "Please set GITHUB_TOKEN and try again."; exit 1; fi
if [ -z "$GOPATH" ]; then echo "Please set GOPATH and try again."; exit 1; fi
if [ -z "$GOBIN" ]; then GOBIN="${GOPATH}/bin"; fi
mkdir -p "${GOBIN}"
echo "Determining platform..."
platform=$(uname | tr '[:upper:]' '[:lower:]')

echo "Finding latest release..."
asset=$(curl --user "x:${GITHUB_TOKEN}" --silent https://api.github.com/repos/owenrumney/squealer/releases/latest | jq -r ".assets[] | select(.name | contains(\"${platform}\")) | .url")
echo "Downloading latest release for your platform..."
curl -s -L -H "Accept: application/octet-stream" --user "x:${GITHUB_TOKEN}" "${asset}" --output ./squealer
echo "Installing squealer..."
chmod +x ./squealer
mv ./squealer "${GOBIN}/"
which squealer &> /dev/null || (echo "Please add ${GOBIN} to your PATH to complete installation!" && exit 1)
echo "Installation complete!"