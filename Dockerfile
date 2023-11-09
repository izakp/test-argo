FROM golang:1.21-alpine3.18

RUN apk add dumb-init git

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build src/server/server.go

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./server"]
