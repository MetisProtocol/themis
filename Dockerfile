# syntax=docker/dockerfile:1
FROM golang:1.22.9 as builder
RUN apt update && apt install -y build-essential git
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build make install

FROM debian:12-slim
COPY --from=builder /go/bin/* /usr/local/bin/
RUN apt-get update && \
    apt-get install -yqq --no-install-recommends sqlite3 curl jq ca-certificates && \
    apt-get clean && rm -rf /var/cache/apt/archives /var/lib/apt/lists/*
EXPOSE 1317 26656 26657
ENTRYPOINT ["themisd", "--home=/var/lib/themis"]
