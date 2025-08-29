CREATE TABLE IF NOT EXISTS orders (
    order_uid        VARCHAR(50) PRIMARY KEY,
    track_number     VARCHAR(50) NOT NULL UNIQUE,
    entry            VARCHAR(50) NOT NULL,
    locale           VARCHAR(10) NOT NULL,
    internal_signature TEXT,
    customer_id      VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey         VARCHAR(50),
    sm_id            INTEGER,
    date_created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    oof_shard        VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS deliveries (
    id        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(50)  NOT NULL,
    phone     VARCHAR(50)  NOT NULL,
    zip       VARCHAR(20)  NOT NULL,
    city      VARCHAR(50)  NOT NULL,
    address   VARCHAR(255) NOT NULL,
    region    VARCHAR(100) NOT NULL,
    email     VARCHAR(100) NOT NULL,
    order_uid VARCHAR(50)  NOT NULL UNIQUE,
    CONSTRAINT fk_deliveries_order
        FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments (
    transaction   VARCHAR(50) PRIMARY KEY,
    request_id    VARCHAR(50),
    currency      CHAR(3)      NOT NULL,
    provider      VARCHAR(50)  NOT NULL,
    amount        DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    payment_dt    TIMESTAMPTZ   NOT NULL,
    bank          VARCHAR(50)   NOT NULL,
    delivery_cost DECIMAL(10,2) NOT NULL CHECK (delivery_cost >= 0),
    goods_total   DECIMAL(10,2) NOT NULL CHECK (goods_total >= 0),
    custom_fee    DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (custom_fee >= 0),
    order_uid     VARCHAR(50)   NOT NULL UNIQUE,
    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id      BIGINT      PRIMARY KEY,    
    price        DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    rid          VARCHAR(50)   NOT NULL,
    name         VARCHAR(100)  NOT NULL,
    sale         SMALLINT      NOT NULL DEFAULT 0 CHECK (sale >= 0),
    size         VARCHAR(10)   NOT NULL,
    total_price  DECIMAL(10,2) NOT NULL CHECK (total_price >= 0),
    nm_id        INTEGER       NOT NULL,
    brand        VARCHAR(100)  NOT NULL,
    status       SMALLINT      NOT NULL,
    track_number VARCHAR(50)   NOT NULL,
    order_uid    VARCHAR(50)   NOT NULL,
    CONSTRAINT fk_items_order
        FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
        ON DELETE CASCADE
);