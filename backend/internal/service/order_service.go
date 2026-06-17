package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var ErrOrderNotFound      = errors.New("pedido não encontrado")
var ErrInvalidOrderStatus = errors.New("status de pedido inválido")

type OrderService struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

func NewOrderService(orderRepo ports.OrderRepository, productRepo ports.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type OrderItemInput struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  float64 `json:"quantity"   validate:"required,gt=0"`
}

type CreateOrderInput struct {
	Items []OrderItemInput `json:"items" validate:"required,min=1,dive"`
	Notes string           `json:"notes"`
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

// ── Operações ────────────────────────────────────────────────────────────────

func (s *OrderService) CreateOrder(ctx context.Context, in CreateOrderInput) (*domain.Order, error) {
	log.Printf("Service - Iniciando CreateOrder com %d itens", len(in.Items))
	order := &domain.Order{
		Status: domain.OrderStatusPending,
		Notes:  in.Notes,
	}

	var total float64

	// Valida produtos e monta os itens com snapshot de preço
	items := make([]domain.OrderItem, len(in.Items))
	for i, itemIn := range in.Items {
		p, err := s.productRepo.FindProductByID(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: buscar produto: %w", err)
		}
		if p == nil || !p.Active {
			return nil, fmt.Errorf("produto id=%d não encontrado ou inativo", itemIn.ProductID)
		}

		items[i] = domain.OrderItem{
			ProductID: p.ID,
			Quantity:  itemIn.Quantity,
			UnitPrice: p.Price, // snapshot do preço atual
		}
		total += p.Price * itemIn.Quantity
	}

	order.Items = items
	order.TotalPrice = total
	log.Printf("Service - Pedido montado: TotalPrice=%f, Items=%d", order.TotalPrice, len(order.Items))

	// Pré-carrega as fichas técnicas antes da transação para evitar context deadline
	productIngredients := make(map[uint][]domain.ProductIngredient)
	for _, itemIn := range in.Items {
		ingredients, err := s.productRepo.GetProductIngredients(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: ficha técnica produto_id=%d: %w", itemIn.ProductID, err)
		}
		productIngredients[itemIn.ProductID] = ingredients
		log.Printf("Service - Produto %d tem %d ingredientes na ficha técnica", itemIn.ProductID, len(ingredients))
	}

	// CreateOrder já executa a baixa de estoque em transação
	// Passamos as fichas técnicas pré-carregadas para evitar chamadas dentro da transação
	if err := s.orderRepo.CreateOrder(ctx, order, productIngredients); err != nil {
		log.Printf("Service - Erro ao criar pedido no repository: %v", err)
		return nil, fmt.Errorf("OrderService.CreateOrder: %w", err)
	}

	log.Printf("Service - Pedido criado com sucesso: ID=%d", order.ID)
	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("OrderService.ListOrders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uint) (*domain.Order, error) {
	order, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.GetOrder: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uint, in UpdateOrderStatusInput) (*domain.Order, error) {
	order, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrderStatus: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	status := domain.OrderStatus(in.Status)
	if err := s.orderRepo.UpdateOrderStatus(ctx, id, status); err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrderStatus: %w", err)
	}
	order.Status = status
	return order, nil
}
