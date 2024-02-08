CREATE TABLE IF NOT EXISTS proxy(
    proxy_id BIGSERIAL PRIMARY KEY,
    protocol VARCHAR(32),
    username VARCHAR(255),
    password VARCHAR(255),
    host VARCHAR(255),
    port BIGINT
);

CREATE TABLE IF NOT EXISTS proxy_occupy (
    key UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    proxy_id BIGINT REFERENCES proxy(proxy_id) ON DELETE CASCADE,
    create_timestamp DOUBLE PRECISION
);