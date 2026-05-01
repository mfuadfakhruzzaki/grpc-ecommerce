FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY proto/ proto/

COPY order-service/go.mod order-service/go.sum ./order-service/
WORKDIR /app/order-service
RUN go mod download
WORKDIR /app
COPY order-service/ order-service/
WORKDIR /app/order-service
RUN go build -o order-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order-service/order-service .
EXPOSE 50053
CMD ["./order-service"]