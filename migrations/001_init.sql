CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    average NUMERIC(20, 8) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_symbol_time ON prices(symbol, recorded_at);

CREATE TABLE source_prices (
    id SERIAL PRIMARY KEY,
    price_id INTEGER REFERENCES prices(id),
    source VARCHAR(50) NOT NULL,
    price NUMERIC(20, 8) NOT NULL
);