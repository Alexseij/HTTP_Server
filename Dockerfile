FROM golang:1.16-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

WORKDIR cmd
WORKDIR server

RUN go build -v -o server

WORKDIR /app

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/cmd/server /app/cmd/server

ENV db_name=FoodFinder
ENV db_user=FoodFinder
ENV db_password=h159357258654
ENV db_host=cluster0.su8jg.mongodb.net
ENV port=8000
ENV host=localhost

EXPOSE 8080

CMD ["/app/cmd/server/server"]