# The instructions for the first stage
FROM node:10-alpine as builder

RUN apk --no-cache add yarn curl bash python3 make g++

# install node-prune (https://github.com/tj/node-prune)
RUN curl -sfL https://install.goreleaser.com/github.com/tj/node-prune.sh | bash -s -- -b /usr/local/bin

WORKDIR /usr/src/app

COPY package.json yarn.lock ./

# install dependencies
RUN yarn --frozen-lockfile

COPY . .

# lint & test
RUN yarn lint & yarn test

# remove development dependencies
RUN npm prune --production

# run node prune
RUN /usr/local/bin/node-prune

# The instructions for second stage
FROM node:10-alpine

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/node_modules ./node_modules
COPY --from=builder /usr/src/app/package.json .

EXPOSE 3000

CMD [ "npm", "start" ]