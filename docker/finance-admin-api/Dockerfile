FROM golang:1.24-alpine AS base
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM base AS dev
WORKDIR /app/finance-admin-api

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base AS build-env
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/finance-admin-api

RUN go build -a -installsuffix cgo -o /go/bin/finance-admin-api

FROM alpine:3
WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3

COPY --from=build-env /go/bin/finance-admin-api finance-admin-api
ENTRYPOINT ["./finance-admin-api"]
