package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type OrderRepository interface {
	// CreateOrder persiste o pedido, seus itens e executa a baixa de estoque
	// de todos os ingredientes em uma única transação atômica.
	// Se qualquer ingrediente ficar com estoque negativo → rollback total.
	// productIngredients é um mapa de product_id -> ingredientes pré-carregados
	CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) error

	FindOrderByID(ctx context.Context, id uint) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error
}
