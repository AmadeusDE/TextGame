# Build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o servergame ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/servergame .
COPY --from=builder /app/data ./data

EXPOSE 8080

CMD ["./servergame"]
