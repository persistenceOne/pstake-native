FROM node:16

WORKDIR /home/explorer
COPY . .

RUN rm -rf node_modules
RUN yarn install --frozen-lockfile

RUN apt-get update -y && apt-get install  ca-certificates jq bash curl -y

EXPOSE 8080
ENTRYPOINT ["yarn","serve"]
