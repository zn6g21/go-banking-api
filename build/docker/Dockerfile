FROM golang:1.25.0-alpine3.22

RUN apk --no-cache add curl

WORKDIR /go/src/web

COPY ./ .

RUN go mod download

ENV GO111MODULE=on

RUN go build -C ./cmd/server/

ENTRYPOINT ./cmd/server/server