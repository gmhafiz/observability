FROM golang:1.22 AS src

WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./


RUN go build -ldflags="-s" -o server .

FROM debian:12

COPY --from=src /go/src/app/server /usr/bin/local/server

EXPOSE 3080

ENTRYPOINT ["/usr/bin/local/server"]
