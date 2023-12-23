-- Таблица заказов
CREATE TABLE IF NOT EXISTS orders (
    order_uid           VARCHAR(255)    UNIQUE PRIMARY KEY,
    track_number        VARCHAR(255),
    entry               VARCHAR(255),
    locale              VARCHAR(255),
    internal_signature  VARCHAR(255),
    customer_id         VARCHAR(255),
    delivery_service    VARCHAR(255),
    shardkey            VARCHAR(255),
    sm_id               INT,
    date_created        TIMESTAMPTZ,
    oof_shard           VARCHAR(255)
);
