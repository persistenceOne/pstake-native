# Use persistenceCore/faucet in order to make this docker image
FROM node:slim as build-env

WORKDIR /home/faucet/
COPY . .

RUN npm install

# Use Alpine and install Node.js which is 50% smaller than the -alpine version of the node
# image (53 MB including the faucet app).
FROM alpine:3.15
# Installs Node.js 16 (https://pkgs.alpinelinux.org/packages?name=nodejs&branch=v3.15)
RUN apk add --update nodejs
RUN apk add --update ca-certificates jq bash curl

WORKDIR /home/faucet

COPY . .

EXPOSE 5001
EXPOSE 5000
