# Build stage
FROM node:20-alpine AS builder

ARG VERSION=v1

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

ENV NODE_ENV=production
ENV APP_VERSION=$VERSION

CMD ["node", "dist/index.js"]