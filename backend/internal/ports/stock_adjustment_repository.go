package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type StockAdjustmentRepository interface {
	// CreateStockAdjustmentPending registra um ajuste de estoque pendente
	// Usado quando um pedido é cancelado para registrar quais ingredientes
	// poderiam ser devolvidos ao estoque após análise manual
	CreateStockAdjustmentPending(ctx context.Context, adjustment *domain.StockAdjustmentPending) error

	// FindPendingByOrderID busca todos os ajustes pendentes para um pedido específico
	FindPendingByOrderID(ctx context.Context, orderID uint) ([]domain.StockAdjustmentPending, error)

	// FindByOrderID busca todos os ajustes (todos os status) para um pedido específico
	FindByOrderID(ctx context.Context, orderID uint) ([]domain.StockAdjustmentPending, error)

	// FindPendingByIngredientID busca ajustes pendentes para um ingrediente específico
	FindPendingByIngredientID(ctx context.Context, ingredientID uint) ([]domain.StockAdjustmentPending, error)

	// ListPending busca todos os ajustes pendentes (para dashboard de aprovação)
	ListPending(ctx context.Context) ([]domain.StockAdjustmentPending, error)

	// UpdateStatus atualiza o status de um ajuste (pending → approved/rejected)
	UpdateStatus(ctx context.Context, id uint, status domain.StockAdjustmentStatus) error

	// FindByID busca um ajuste por ID
	FindByID(ctx context.Context, id uint) (*domain.StockAdjustmentPending, error)
}
