-- +goose Up
-- +goose StatementBegin

-- Adiciona constraint de unicidade para prevenir duplicatas
-- Um pedido não pode ter múltiplos ajustes pendentes para o mesmo ingrediente
CREATE UNIQUE INDEX IF NOT EXISTS uk_stock_adjustments_order_ingredient_pending
ON stock_adjustments_pending(order_id, ingredient_id)
WHERE status = 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS uk_stock_adjustments_order_ingredient_pending;

-- +goose StatementEnd
