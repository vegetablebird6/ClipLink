# syntax=docker/dockerfile:1

FROM node:22-alpine AS web-builder
WORKDIR /src/web

COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

FROM golang:1.25-alpine AS go-builder
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=web-builder /src/web/dist ./internal/static/dist

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/cliplink \
    ./cmd/main.go

FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S cliplink \
    && adduser -S -G cliplink -h /home/cliplink cliplink \
    && mkdir -p /home/cliplink/.cliplink /app \
    && chown -R cliplink:cliplink /home/cliplink /app

WORKDIR /app

COPY --from=go-builder /out/cliplink /app/cliplink
COPY config.example.yml /app/config.example.yml

ENV GIN_MODE=release \
    HOME=/home/cliplink

EXPOSE 8080
VOLUME ["/home/cliplink/.cliplink"]

USER cliplink

ENTRYPOINT ["/app/cliplink"]
