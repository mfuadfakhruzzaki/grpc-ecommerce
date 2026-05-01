FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY proto/ proto/

COPY gateway/go.mod gateway/go.sum ./gateway/
WORKDIR /app/gateway
RUN go mod download
WORKDIR /app
COPY gateway/ gateway/
WORKDIR /app/gateway
RUN go build -o gateway .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway/gateway .
EXPOSE 8080
CMD ["./gateway"]