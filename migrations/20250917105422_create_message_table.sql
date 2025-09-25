-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL CHECK (type IN ('доход', 'расход')),
    amount INT NOT NULL CHECK (amount >= 0),
    date DATE NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION percentile_cont_window(vals double precision[], pct double precision)
RETURNS double precision
LANGUAGE sql
AS $$
    WITH unnest_values AS (
        SELECT unnest(vals) AS val
    )
    SELECT percentile_cont(pct) WITHIN GROUP (ORDER BY val)
    FROM unnest_values;
$$;

-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_items_date ON items (date);
CREATE INDEX idx_items_type ON items (type);
CREATE INDEX idx_items_category ON items (category);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS images;
