CREATE TABLE IF NOT EXISTS proxies(
    id serial PRIMARY KEY,
    protocol VARCHAR(32),
    username VARCHAR(255),
    password VARCHAR(255),
    host VARCHAR(255),
    port INTEGER,
    clients_count INTEGER DEFAULT(0)
);