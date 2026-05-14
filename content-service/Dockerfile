
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o content-service ./cmd/content-service

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/content-service .

EXPOSE 50052 

CMD ["./content-service"]