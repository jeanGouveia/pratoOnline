package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var ErrStockAdjustmentNotFound = errors.New("ajuste de estoque não encontrado")

type StockAdjustmentService struct {
	stockAdjustmentRepo ports.StockAdjustmentRepository
	productRepo         ports.ProductRepository
}

func NewStockAdjustmentService(
	stockAdjustmentRepo ports.StockAdjustmentRepository,
	productRepo ports.ProductRepository,
) *StockAdjustmentService {
	return &StockAdjustmentService{
		stockAdjustmentRepo: stockAdjustmentRepo,
		productRepo:         productRepo,
	}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type CreateStockAdjustmentInput struct {
	OrderID      uint    `json:"order_id" validate:"required"`
	IngredientID uint    `json:"ingredient_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	OrderStatus  string  `json:"order_status" validate:"required"`
}

type UpdateStockAdjustmentStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending approved rejected"`
}

// ── Operações ────────────────────────────────────────────────────────────────

// RegisterStockAdjustmentForOrder registra ajustes de estoque pendentes para um pedido cancelado
// Este método NÃO devolve estoque automaticamente, apenas registra para análise manual
func (s *StockAdjustmentService) RegisterStockAdjustmentForOrder(
	ctx context.Context,
	orderID uint,
	orderStatus domain.OrderStatus,
	productIngredients map[uint][]domain.ProductIngredient,
	orderItems []domain.OrderItem,
) error {
	// Validar se já existem ajustes pendentes para este pedido (prevenção de duplicatas)
	existing, err := s.stockAdjustmentRepo.FindPendingByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("verificar ajustes existentes: %w", err)
	}
	if len(existing) > 0 {
		return fmt.Errorf("já existem %d ajustes pendentes para o pedido %d", len(existing), orderID)
	}

	// Para cada item do pedido
	for _, item := range orderItems {
		ingredients, ok := productIngredients[item.ProductID]
		if !ok || len(ingredients) == 0 {
			// Produto simples ou sem ficha técnica: nada a registrar
			continue
		}

		// Para cada ingrediente da ficha técnica
		for _, pi := range ingredients {
			// Calcular quantidade consumida
			consumedQuantity := pi.Quantity * item.Quantity

			// Registrar ajuste pendente
			adjustment := &domain.StockAdjustmentPending{
				OrderID:      orderID,
				IngredientID: pi.IngredientID,
				Quantity:     consumedQuantity,
				OrderStatus:  string(orderStatus),
				Status:       domain.StockAdjustmentStatusPending,
			}

			if err := s.stockAdjustmentRepo.CreateStockAdjustmentPending(ctx, adjustment); err != nil {
				return fmt.Errorf("RegisterStockAdjustmentForOrder: %w", err)
			}
		}
	}

	return nil
}

// ListPendingAdjustments lista todos os ajustes pendentes para aprovação
func (s *StockAdjustmentService) ListPendingAdjustments(
	ctx context.Context,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.ListPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListPendingAdjustments: %w", err)
	}
	return adjustments, nil
}

// GetAdjustmentsByOrder lista todos os ajustes (todos os status) para um pedido
func (s *StockAdjustmentService) GetAdjustmentsByOrder(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("GetAdjustmentsByOrder: %w", err)
	}
	return adjustments, nil
}

// GetPendingAdjustmentsByOrder lista apenas os ajustes pendentes para um pedido
func (s *StockAdjustmentService) GetPendingAdjustmentsByOrder(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindPendingByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("GetPendingAdjustmentsByOrder: %w", err)
	}
	return adjustments, nil
}

// GetPendingAdjustmentsByIngredient lista ajustes pendentes para um ingrediente específico
func (s *StockAdjustmentService) GetPendingAdjustmentsByIngredient(
	ctx context.Context, ingredientID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindPendingByIngredientID(ctx, ingredientID)
	if err != nil {
		return nil, fmt.Errorf("GetPendingAdjustmentsByIngredient: %w", err)
	}
	return adjustments, nil
}

// GetAdjustmentByID busca um ajuste específico por ID
func (s *StockAdjustmentService) GetAdjustmentByID(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	adjustment, err := s.stockAdjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetAdjustmentByID: %w", err)
	}
	if adjustment == nil {
		return nil, ErrStockAdjustmentNotFound
	}
	return adjustment, nil
}

// ApproveAdjustment aprova um ajuste de estoque pendente
// NOTA: Este método apenas marca como aprovado. A devolução efetiva do estoque
// deve ser implementada separadamente com validações específicas de negócio
func (s *StockAdjustmentService) ApproveAdjustment(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	adjustment, err := s.GetAdjustmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if adjustment == nil {
		return nil, ErrStockAdjustmentNotFound
	}

	if err := s.stockAdjustmentRepo.UpdateStatus(ctx, id, domain.StockAdjustmentStatusApproved); err != nil {
		return nil, fmt.Errorf("ApproveAdjustment: %w", err)
	}

	adjustment.Status = domain.StockAdjustmentStatusApproved
	return adjustment, nil
}

// RejectAdjustment rejeita um ajuste de estoque pendente
func (s *StockAdjustmentService) RejectAdjustment(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	adjustment, err := s.GetAdjustmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if adjustment == nil {
		return nil, ErrStockAdjustmentNotFound
	}

	if err := s.stockAdjustmentRepo.UpdateStatus(ctx, id, domain.StockAdjustmentStatusRejected); err != nil {
		return nil, fmt.Errorf("RejectAdjustment: %w", err)
	}

	adjustment.Status = domain.StockAdjustmentStatusRejected
	return adjustment, nil
}
