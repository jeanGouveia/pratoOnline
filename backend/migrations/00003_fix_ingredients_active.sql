-- +goose Up

ALTER TABLE ingredients
ADD COLUMN active INTEGER NOT NULL DEFAULT 1;

-- +goose Down

SELECT 1;