FROM golang:1.17-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /usr/local/share/app

# Add source files
COPY docker/exposes .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && go build -o exposer ./...

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates jq bash
WORKDIR /usr/local/share/app

# Copy over binaries from the build-env
COPY --from=build-env /usr/local/share/app/exposer /usr/bin/exposer

EXPOSE 8080
