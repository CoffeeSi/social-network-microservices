FROM golang:1.26.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o chat-service ./cmd/chat-service/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/chat-service .

EXPOSE 50053

CMD ["./chat-service"]