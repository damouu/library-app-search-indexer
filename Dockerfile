# Build stage
FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o search-indexer ./cmd/search-indexer


# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/search-indexer .

CMD ["./search-indexer"]