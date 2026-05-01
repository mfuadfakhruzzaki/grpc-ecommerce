FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY proto/ proto/

COPY product-service/go.mod product-service/go.sum ./product-service/
WORKDIR /app/product-service
RUN go mod download
WORKDIR /app
COPY product-service/ product-service/
WORKDIR /app/product-service
RUN go build -o product-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/product-service/product-service .
EXPOSE 50052
CMD ["./product-service"]