FROM golang:1.16.3-alpine  AS build-env

ARG squealer_version=0.0.0

COPY . /src
WORKDIR /src
ENV CGO_ENABLED=0
RUN go build \
  -a \
  -ldflags "-X github.com/owenrumney/squealer/version.Version=${squealer_version} -s -w -extldflags '-static'" \
  -mod=vendor \
  ./cmd/squealer


FROM alpine

# use a non-privileged user
USER nobody

# work somewhere where we can write
COPY --from=build-env /src/squealer /usr/bin/squealer

# set the default entrypoint -- when this container is run, use this command
ENTRYPOINT [ "sqealer" ]

# as we specified an entrypoint, this is appended as an argument (i.e., `sqealer --help`)
CMD [ "--help" ]