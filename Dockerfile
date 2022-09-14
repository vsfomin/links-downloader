# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.1-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /links-downloader cmd/main.go 

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /links-downloader /links-downloader

ENTRYPOINT ["/links-downloader"]