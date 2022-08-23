FROM golang:1.18-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go env -w GOFLAGS=-buildvcs=false

CMD ["go", "run", "server.go"]