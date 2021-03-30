# The instructions for the first stage
FROM node:10-alpine as builder

RUN apk --no-cache add curl bash python3 make g++

# install node-prune (https://github.com/tj/node-prune)
RUN curl -sfL https://install.goreleaser.com/github.com/tj/node-prune.sh | bash -s -- -b /usr/local/bin

COPY package*.json ./
RUN npm install

# remove development dependencies
RUN npm prune --production

# run node prune
RUN /usr/local/bin/node-prune

# The instructions for second stage
FROM node:10-alpine

WORKDIR /usr/src/app
COPY --from=builder node_modules node_modules

EXPOSE 3000

CMD [ "npm", "run", "start-prod" ]