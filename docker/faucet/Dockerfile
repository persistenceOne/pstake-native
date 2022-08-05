FROM node:slim as build-env

WORKDIR /faucet/
COPY . .

RUN npm ci
RUN npx esbuild server.js --bundle --platform=node --outfile=index.js

FROM node:slim

WORKDIR /app/

COPY --from=build-env /faucet/index.js /app/index.js

CMD node /app/index.js

EXPOSE 5000
