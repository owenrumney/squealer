#!/usr/bin/env bash

set -e

IMAGE=owenrumney/squealer
docker login -u $DOCKER_USERNAME --password $DOCKER_PASSWORD

echo "building ${IMAGE}..."
docker build --build-arg squealer_version=${TRAVIS_TAG} -t ${IMAGE} .

echo "publishing ${IMAGE}..."
# push the patch tag - eg; v0.36.15
docker tag ${IMAGE} ${IMAGE}:${TRAVIS_TAG}
docker push ${IMAGE}:${TRAVIS_TAG}

# push the minor tag - eg; v0.36
docker tag ${IMAGE} ${IMAGE}:${TRAVIS_TAG%.*}
docker push ${IMAGE}:${TRAVIS_TAG%.*}

# push the latest tag
docker tag ${IMAGE} ${IMAGE}:latest
docker push ${IMAGE}:latest