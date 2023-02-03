FROM golang:1.19-alpine as build

WORKDIR /src/app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build

FROM alpine:edge

COPY --from=build /src/app/interchain-queries /usr/bin/interchain-queries

RUN apk add --update ca-certificates jq bash curl
WORKDIR /root

RUN adduser --system --home /icq --disabled-password --disabled-login icq -U 1000

USER icq
