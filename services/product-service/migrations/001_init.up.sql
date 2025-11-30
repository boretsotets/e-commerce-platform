CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    price REAL,
    stock INT
);