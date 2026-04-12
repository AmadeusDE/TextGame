# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o servergame ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/servergame .
COPY --from=builder /app/data ./data
COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./servergame"]
