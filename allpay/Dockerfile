FROM golang:1.19 AS build
WORKDIR /go/src
COPY allpay/go ./go
COPY allpay/main.go .
COPY allpay/go.sum .
COPY allpay/go.mod .

ENV CGO_ENABLED=0

RUN go build -o openapi .

FROM scratch AS runtime
COPY --from=build /go/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
