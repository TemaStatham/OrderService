-- Таблица товаров
CREATE TABLE IF NOT EXISTS order_items (
    id                      SERIAL PRIMARY KEY,
    chrt_id                 INT,
    order_uid               VARCHAR(255),
    track_number            VARCHAR(255),
    price                   INT,
    rid                     VARCHAR(255),
    name                    VARCHAR(255),
    sale                    INT,
    size                    VARCHAR(255),
    total_price             INT,
    nm_id                   INT,
    brand                   VARCHAR(255),
    status                  INT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);
