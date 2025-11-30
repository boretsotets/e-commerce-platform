CREATE TABLE order_dbs (
    order_id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    "status" TEXT NOT NULL,
    shipping_address TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);