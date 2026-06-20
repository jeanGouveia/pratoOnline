package domain

import "time"

// StockAdjustmentStatus representa o status de um ajuste de estoque pendente
type StockAdjustmentStatus string

const (
	StockAdjustmentStatusPending  StockAdjustmentStatus = "pending"
	StockAdjustmentStatusApproved StockAdjustmentStatus = "approved"
	StockAdjustmentStatusRejected StockAdjustmentStatus = "rejected"
)

// StockAdjustmentPending registra ajustes de estoque que precisam de aprovação manual
// ou análise antes de serem aplicados. Usado para estornos de estoque por cancelamento
// de pedidos, permitindo auditoria e validação humana.
type StockAdjustmentPending struct {
	ID           uint
	OrderID      uint
	IngredientID uint
	Quantity     float64 // Quantidade que poderia ser devolvida ao estoque
	OrderStatus  string  // Status do pedido no momento do cancelamento (para contexto)
	Status       StockAdjustmentStatus
	CreatedAt    time.Time
	ProcessedAt  *time.Time // Quando foi aprovado/rejeitado (null se pending)
}
