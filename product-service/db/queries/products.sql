-- name: CreateProduct :one
INSERT INTO products (name, description, price, category_id, stock_qty)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND is_active = true LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE is_active = true;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products SET is_active = false WHERE id = $1;

-- name: GetStock :one
SELECT stock_qty FROM products WHERE id = $1;

-- name: DeductStock :one
UPDATE products
SET stock_qty = stock_qty - $2, updated_at = NOW()
WHERE id = $1 AND stock_qty >= $2
RETURNING stock_qty;