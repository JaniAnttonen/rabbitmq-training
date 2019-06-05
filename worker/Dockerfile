# Build
FROM golang:alpine AS build

RUN apk update && apk add git

RUN adduser -D -g '' workeruser

RUN mkdir -p /go/src/service-worker
COPY . /go/src/service-worker

ENV GOPATH /go:$GOPATH
ENV PATH /go/bin:$PATH

RUN go get github.com/streadway/amqp

WORKDIR /go/src/service-worker

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/service-worker

# Release
FROM scratch

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/service-worker /go/bin/service-worker

USER workeruser
