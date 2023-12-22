FROM golang:1.21-bookworm as builder

WORKDIR /app

COPY . ./
RUN go mod tidy

# Build the binary.
RUN go build -v -o server

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

CMD ["/app/server"]
