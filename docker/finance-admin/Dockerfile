FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS base

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    app-user

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM node:24-alpine3.23@sha256:0388af2af070cd4736a1567cfed02469ba117848845b4165d87a333edb53d2ca AS asset-env
WORKDIR /app

RUN mkdir -p web/static

RUN apk add git

COPY finance-admin/web/assets web/assets
COPY finance-admin/package.json .
COPY finance-admin/package-lock.json .

RUN npm ci --ignore-scripts
RUN npm run build

FROM base AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck

RUN go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM base AS dev
WORKDIR /app/finance-admin

COPY --from=asset-env /app/web/static web/static
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=5 CMD [ "/app/finance-admin/healthcheck" ]

ENTRYPOINT ["air"]

FROM base AS app-build
WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY healthcheck healthcheck
COPY shared shared
COPY finance-admin finance-admin
COPY apierror apierror

RUN go mod download

WORKDIR /app/finance-admin

RUN go build -a -installsuffix cgo -o /go/bin/finance-admin

FROM scratch 

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=app-build /go/bin/finance-admin finance-admin
COPY --from=app-build /app/finance-admin/web/template web/template
COPY --from=asset-env /app/web/static web/static
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck

USER app-user:app-user

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=5 CMD [ "/healthcheck" ]

ENTRYPOINT ["./finance-admin"]
