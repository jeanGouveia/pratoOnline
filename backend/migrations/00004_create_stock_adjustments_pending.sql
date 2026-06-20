-- +goose Up
-- +goose StatementBegin

-- Tabela para registrar ajustes de estoque pendentes de aprovação
-- Usada para estornos de estoque por cancelamento de pedidos
CREATE TABLE IF NOT EXISTS stock_adjustments_pending (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id      INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    quantity      REAL NOT NULL,
    order_status  TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'pending',
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    processed_at  DATETIME
);

-- Índices para consultas frequentes
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_order_id ON stock_adjustments_pending(order_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_ingredient_id ON stock_adjustments_pending(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_status ON stock_adjustments_pending(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_stock_adjustments_pending_status;
DROP INDEX IF EXISTS idx_stock_adjustments_pending_ingredient_id;
DROP INDEX IF EXISTS idx_stock_adjustments_pending_order_id;
DROP TABLE IF EXISTS stock_adjustments_pending;

-- +goose StatementEnd
