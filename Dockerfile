FROM golang:1.16-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

WORKDIR /app/cmd/server

RUN go build

WORKDIR /app

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/cmd/server cmd/server
COPY --from=builder /app/config config

WORKDIR /app/cmd/server

RUN ["./server"]