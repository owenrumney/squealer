FROM alpine:3.21.2

# use a non-privileged user
USER nobody

# work somewhere where we can write
COPY squealer /usr/bin/squealer

# set the default entrypoint -- when this container is run, use this command
ENTRYPOINT [ "sqealer" ]

# as we specified an entrypoint, this is appended as an argument (i.e., `sqealer --help`)
CMD [ "--help" ]


