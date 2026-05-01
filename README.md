# gRPC E-Commerce Microservices

Sistem belanja online berbasis microservices dengan gRPC, API Gateway, dan Load Balancer.

## Arsitektur

```
Client (REST)
     │
     ▼
┌─────────────┐
│ Nginx :80   │  ← Load Balancer
└─────────────┘
     │
     ▼
┌─────────────┐
│ Gateway     │  ← HTTP/REST → gRPC, JWT Auth, Rate Limit
│ :8080       │
└─────────────┘
     │
     ├──────────────────┬──────────────────┐
     ▼                  ▼                  ▼
┌──────────┐     ┌──────────────┐    ┌──────────────┐
│  User    │     │   Product    │    │    Order     │
│ Service  │     │   Service    │    │   Service    │
│  :50051  │     │   :50052     │    │   :50053     │
└──────────┘     └──────────────┘    └──────────────┘
     │                  │                  │
     ▼                  ▼                  ▼
┌──────────┐     ┌──────────────┐    ┌──────────────┐
│ user_db  │     │  product_db  │    │   order_db   │
│ :5432    │     │  :5433       │    │   :5434      │
└──────────┘     └──────────────┘    └──────────────┘
```

## Tech Stack

| Layer         | Teknologi                  |
| ------------- | -------------------------- |
| Language      | Go 1.25+                   |
| Inter-service | gRPC + Protocol Buffers v3 |
| REST Gateway  | grpc-gateway v2            |
| Load Balancer | Nginx                      |
| Database      | PostgreSQL 16              |
| Query         | sqlc + golang-migrate      |
| Auth          | JWT (golang-jwt/jwt)       |
| Container     | Docker + Docker Compose    |

## Cara Menjalankan

### Prasyarat

- Docker & Docker Compose
- Go 1.25+
- protoc + plugins
- sqlc
- golang-migrate

### Run lokal

```bash
# Clone repo
git clone https://github.com/mfuadfakhruzzaki/grpc-ecommerce
cd grpc-ecommerce

# Jalankan semua service
docker compose up --build
```

Semua service akan berjalan otomatis termasuk database dan migrasi.

### Endpoints

Base URL: `http://localhost:8080`

#### Auth

```bash
# Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","full_name":"John Doe"}'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

#### Users (butuh JWT)

```bash
# Get profile
curl http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer <token>"

# Update profile
curl -X PATCH http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Jane Doe"}'
```

#### Products (butuh JWT)

```bash
# List products
curl http://localhost:8080/v1/products \
  -H "Authorization: Bearer <token>"

# Create product
curl -X POST http://localhost:8080/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Produk A","description":"Deskripsi","price":50000,"stock_qty":100}'

# Get product
curl http://localhost:8080/v1/products/<id> \
  -H "Authorization: Bearer <token>"

# Update product
curl -X PATCH http://localhost:8080/v1/products/<id> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Produk B","price":60000}'

# Delete product
curl -X DELETE http://localhost:8080/v1/products/<id> \
  -H "Authorization: Bearer <token>"
```

#### Orders (butuh JWT)

```bash
# Create order
curl -X POST http://localhost:8080/v1/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":"<id>","quantity":2}]}'

# Get order
curl http://localhost:8080/v1/orders/<id> \
  -H "Authorization: Bearer <token>"

# List orders
curl http://localhost:8080/v1/orders \
  -H "Authorization: Bearer <token>"

# Update order status
curl -X PATCH http://localhost:8080/v1/orders/<id>/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

## Struktur Proyek

```
grpc-ecommerce/
├── proto/                  # Proto definitions
│   ├── user/v1/
│   ├── product/v1/
│   └── order/v1/
├── gateway/                # API Gateway (HTTP → gRPC)
│   ├── middleware/         # JWT auth, rate limit
│   └── main.go
├── user-service/
│   ├── internal/
│   │   ├── handler/        # gRPC handlers
│   │   ├── repository/     # DB queries (sqlc)
│   │   └── service/        # Business logic
│   ├── db/migrations/
│   └── main.go
├── product-service/        # Struktur sama
├── order-service/          # Struktur sama
├── nginx/nginx.conf        # Load balancer config
├── docker-compose.yml
├── Makefile
└── README.md
```

## Fitur

- ✅ Registrasi & login dengan JWT
- ✅ CRUD produk dengan manajemen stok
- ✅ Pembuatan order dengan inter-service gRPC call
- ✅ Tracking status order
- ✅ API Gateway dengan JWT validation
- ✅ Load balancer dengan Nginx
- ✅ Database-per-service pattern
- ✅ Single command deployment: `docker compose up`

## Author

**mfuadfakhruzzaki** — [github.com/mfuadfakhruzzaki](https://github.com/mfuadfakhruzzaki)
