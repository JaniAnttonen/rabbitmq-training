# Build
FROM golang:alpine AS build

RUN apk update && apk add git

RUN adduser -D -g '' apiuser

RUN mkdir -p /go/src/service-ingester
COPY . /go/src/service-ingester

ENV GOPATH /go:$GOPATH
ENV PATH /go/bin:$PATH

RUN go get github.com/streadway/amqp github.com/vincentLiuxiang/lu github.com/valyala/fasthttp

WORKDIR /go/src/service-ingester

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/service-ingester

# Release
FROM scratch

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/service-ingester /go/bin/service-ingester

USER apiuser
