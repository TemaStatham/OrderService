-- Таблица оплаты
CREATE TABLE IF NOT EXISTS order_payment (
    id                          SERIAL PRIMARY KEY,
    order_uid                   VARCHAR(255),
    transaction                 VARCHAR(255),
    request_id                  VARCHAR(255),
    currency                    VARCHAR(255),
    provider                    VARCHAR(255),
    amount                      INT,
    payment_dt                  INT,
    bank                        VARCHAR(255),
    delivery_cost               INT,
    goods_total                 INT,
    custom_fee                  INT,
    FOREIGN KEY (order_uid)     REFERENCES orders(order_uid)
);