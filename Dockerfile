FROM golang:1.18-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir -p /app
WORKDIR /app

COPY . .

RUN go mod download

RUN go env -w GOFLAGS=-buildvcs=false

CMD ["go", "run", "server.go"]