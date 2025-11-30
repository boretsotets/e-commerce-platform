CREATE TABLE order_item_dbs (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGSERIAL NOT NULL REFERENCES order_dbs(order_id),
    product_id BIGINT NOT NULL,
    count INT NOT NULL
);
