FROM golang:1.17-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/peristenceOne/pStake-native/orchestrator

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && cd orchestrator && \
  make

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates jq bash curl
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/orchestrator /usr/bin/orchestrator


