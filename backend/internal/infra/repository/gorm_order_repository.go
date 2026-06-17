package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type gormOrder struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Status     string `gorm:"not null;default:'pending'"`
	TotalPrice float64
	Notes      string
	CreatedAt  int64  `gorm:"autoCreateTime"`
	UpdatedAt  int64  `gorm:"autoUpdateTime"`
}

func (gormOrder) TableName() string { return "orders" }

type gormOrderItem struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null"`
	Quantity  float64 `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
}

func (gormOrderItem) TableName() string { return "order_items" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.OrderRepository = (*GormOrderRepository)(nil)

type GormOrderRepository struct {
	db          *gorm.DB
	productRepo ports.ProductRepository
}

func NewGormOrderRepository(db *gorm.DB, productRepo ports.ProductRepository) *GormOrderRepository {
	return &GormOrderRepository{db: db, productRepo: productRepo}
}

// CreateOrder é a operação crítica: persiste pedido + itens + baixa de estoque
// em uma única transação. Qualquer falha reverte tudo.
func (r *GormOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) error {
	log.Printf("Repository - Iniciando CreateOrder com %d itens", len(order.Items))
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Persiste o pedido
		gOrder := gormOrder{
			Status:     string(order.Status),
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
		}
		if err := tx.Create(&gOrder).Error; err != nil {
			log.Printf("Repository - Erro ao criar pedido: %v", err)
			return fmt.Errorf("CreateOrder: criar pedido: %w", err)
		}
		order.ID = gOrder.ID
		order.CreatedAt = time.Unix(gOrder.CreatedAt, 0)
		log.Printf("Repository - Pedido criado no banco: ID=%d", order.ID)

		// 2. Para cada item do pedido
		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID
			log.Printf("Repository - Processando item: ProductID=%d, Quantity=%f", item.ProductID, item.Quantity)

			// 2a. Persiste o item
			gItem := gormOrderItem{
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			}
			if err := tx.Create(&gItem).Error; err != nil {
				log.Printf("Repository - Erro ao criar item: %v", err)
				return fmt.Errorf("CreateOrder: criar item produto_id=%d: %w", item.ProductID, err)
			}
			item.ID = gItem.ID
			log.Printf("Repository - Item criado: ID=%d", item.ID)

			// 2b. Usa os ingredientes pré-carregados (evita chamada dentro da transação)
			ingredients, ok := productIngredients[item.ProductID]
			if !ok {
				log.Printf("Repository - Ingredientes não pré-carregados para produto_id=%d", item.ProductID)
				return fmt.Errorf("CreateOrder: ingredientes não pré-carregados para produto_id=%d", item.ProductID)
			}

			if len(ingredients) == 0 {
				// Produto simples sem ficha técnica: não há ingredientes para baixar
				log.Printf("Repository - Produto %d não tem ficha técnica (produto simples)", item.ProductID)
				continue
			}

			// 2c. Para cada ingrediente da ficha técnica, dá baixa proporcional
			// Consumo = quantidade_na_ficha × quantidade_vendida
			for _, pi := range ingredients {
				consumo := pi.Quantity * item.Quantity
				log.Printf("Repository - Baixando estoque: IngredientID=%d, Consumo=%f", pi.IngredientID, consumo)
				if err := r.productRepo.DecreaseIngredientStock(ctx, pi.IngredientID, consumo, tx); err != nil {
					log.Printf("Repository - Erro ao baixar estoque: %v", err)
					// Erro retorna com nome do ingrediente já embutido
					return fmt.Errorf("CreateOrder: baixa estoque: %w", err)
				}
			}
		}

		log.Printf("Repository - Transação commitada com sucesso")
		return nil // commit
	})
}

func (r *GormOrderRepository) FindOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	var gOrder gormOrder
	err := r.db.WithContext(ctx).First(&gOrder, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindOrderByID: %w", err)
	}

	var gItems []gormOrderItem
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).Find(&gItems).Error; err != nil {
		return nil, fmt.Errorf("FindOrderByID items: %w", err)
	}

	order := orderToDomain(&gOrder)
	order.Items = make([]domain.OrderItem, len(gItems))
	for i, gi := range gItems {
		order.Items[i] = domain.OrderItem{
			ID: gi.ID, OrderID: gi.OrderID,
			ProductID: gi.ProductID, Quantity: gi.Quantity, UnitPrice: gi.UnitPrice,
		}
	}
	return order, nil
}

func (r *GormOrderRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var gOrders []gormOrder
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&gOrders).Error; err != nil {
		return nil, fmt.Errorf("ListOrders: %w", err)
	}
	out := make([]domain.Order, len(gOrders))
	for i, g := range gOrders {
		out[i] = *orderToDomain(&g)
	}
	return out, nil
}

func (r *GormOrderRepository) UpdateOrderStatus(
	ctx context.Context, id uint, status domain.OrderStatus,
) error {
	if err := r.db.WithContext(ctx).Model(&gormOrder{}).
		Where("id = ?", id).Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("UpdateOrderStatus: %w", err)
	}
	return nil
}

func orderToDomain(g *gormOrder) *domain.Order {
	return &domain.Order{
		ID: g.ID, Status: domain.OrderStatus(g.Status),
		TotalPrice: g.TotalPrice, Notes: g.Notes,
		CreatedAt: time.Unix(g.CreatedAt, 0), UpdatedAt: time.Unix(g.UpdatedAt, 0),
	}
}
