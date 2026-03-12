# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o backfeedr ./cmd/backfeedr

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/backfeedr .
EXPOSE 8080
VOLUME ["/data"]
ENV BACKFEEDR_DB_PATH=/data/backfeedr.db
ENTRYPOINT ["./backfeedr"]
