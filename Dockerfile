FROM alpine:3.23.3

# Set the working directory
WORKDIR /usr/bin

# Copy the binary to the working directory
COPY squealer /usr/bin/squealer

# Set the default entrypoint
ENTRYPOINT [ "squealer" ]