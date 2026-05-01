FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy root go.mod dulu (untuk proto)
COPY go.mod go.sum ./
COPY proto/ proto/

# Copy user-service
COPY user-service/go.mod user-service/go.sum ./user-service/
WORKDIR /app/user-service
RUN go mod download
WORKDIR /app
COPY user-service/ user-service/
WORKDIR /app/user-service
RUN go build -o user-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service/user-service .
EXPOSE 50051
CMD ["./user-service"]