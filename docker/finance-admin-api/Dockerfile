FROM golang:1.24-alpine AS base

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

FROM base AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck

RUN go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM base AS dev
WORKDIR /app/finance-admin-api

COPY --from=healthcheck-build /go/bin/healthcheck healthcheck

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=5 CMD [ "/app/finance-admin-api/healthcheck" ]

ENTRYPOINT ["air"]

FROM base AS app-build
WORKDIR /app

COPY go.mod go.sum ./

COPY healthcheck healthcheck
COPY shared shared
COPY finance-admin-api finance-admin-api
COPY apierror apierror

RUN go mod download

WORKDIR /app/finance-admin-api

RUN go build -a -installsuffix cgo -o /go/bin/finance-admin-api

FROM scratch 

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=app-build /go/bin/finance-admin-api finance-admin-api
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck

USER app-user:app-user

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=5 CMD [ "/healthcheck" ]

ENTRYPOINT ["./finance-admin-api"]
