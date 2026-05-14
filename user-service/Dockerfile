FROM golang:1.26.2-alpine AS builder


RUN apk add --no-cache git ca-certificates

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 50050
CMD ["./main"]