CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE order_status AS ENUM (
    'pending', 'confirmed', 'shipped', 'delivered', 'cancelled'
);

CREATE TABLE orders (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL,
    status       order_status NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12,2) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id       UUID REFERENCES orders(id),
    product_id     UUID NOT NULL,
    quantity       INTEGER NOT NULL,
    price_at_order NUMERIC(12,2) NOT NULL
);