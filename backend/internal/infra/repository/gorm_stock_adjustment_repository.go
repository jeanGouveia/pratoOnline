package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type gormStockAdjustmentPending struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	OrderID      uint    `gorm:"not null;index"`
	IngredientID uint    `gorm:"not null;index"`
	Quantity     float64 `gorm:"not null"`
	OrderStatus  string  `gorm:"not null"`
	Status       string  `gorm:"not null;default:'pending';index"`
	CreatedAt    int64   `gorm:"autoCreateTime"`
	ProcessedAt  *int64  `gorm:"index"`
}

func (gormStockAdjustmentPending) TableName() string { return "stock_adjustments_pending" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.StockAdjustmentRepository = (*GormStockAdjustmentRepository)(nil)

type GormStockAdjustmentRepository struct {
	db *gorm.DB
}

func NewGormStockAdjustmentRepository(db *gorm.DB) *GormStockAdjustmentRepository {
	return &GormStockAdjustmentRepository{db: db}
}

func (r *GormStockAdjustmentRepository) CreateStockAdjustmentPending(
	ctx context.Context, adjustment *domain.StockAdjustmentPending,
) error {
	gAdjustment := gormStockAdjustmentPending{
		OrderID:      adjustment.OrderID,
		IngredientID: adjustment.IngredientID,
		Quantity:     adjustment.Quantity,
		OrderStatus:  adjustment.OrderStatus,
		Status:       string(adjustment.Status),
	}
	if err := r.db.WithContext(ctx).Create(&gAdjustment).Error; err != nil {
		return fmt.Errorf("CreateStockAdjustmentPending: %w", err)
	}
	adjustment.ID = gAdjustment.ID
	adjustment.CreatedAt = time.Unix(gAdjustment.CreatedAt, 0)
	return nil
}

func (r *GormStockAdjustmentRepository) FindPendingByOrderID(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []gormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("order_id = ? AND status = ?", orderID, domain.StockAdjustmentStatusPending).
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindPendingByOrderID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) FindByOrderID(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []gormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindByOrderID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) FindPendingByIngredientID(
	ctx context.Context, ingredientID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []gormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("ingredient_id = ? AND status = ?", ingredientID, domain.StockAdjustmentStatusPending).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindPendingByIngredientID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) ListPending(
	ctx context.Context,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []gormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("status = ?", domain.StockAdjustmentStatusPending).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("ListPending: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) UpdateStatus(
	ctx context.Context, id uint, status domain.StockAdjustmentStatus,
) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status": string(status),
	}
	if status != domain.StockAdjustmentStatusPending {
		updates["processed_at"] = now
	}
	if err := r.db.WithContext(ctx).
		Model(&gormStockAdjustmentPending{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}
	return nil
}

func (r *GormStockAdjustmentRepository) FindByID(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	var gAdjustment gormStockAdjustmentPending
	if err := r.db.WithContext(ctx).First(&gAdjustment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return r.mapToDomain(&gAdjustment), nil
}

// ─── Mappers ──────────────────────────────────────────────────────────────────

func (r *GormStockAdjustmentRepository) mapToDomain(g *gormStockAdjustmentPending) *domain.StockAdjustmentPending {
	adjustment := &domain.StockAdjustmentPending{
		ID:           g.ID,
		OrderID:      g.OrderID,
		IngredientID: g.IngredientID,
		Quantity:     g.Quantity,
		OrderStatus:  g.OrderStatus,
		Status:       domain.StockAdjustmentStatus(g.Status),
		CreatedAt:    time.Unix(g.CreatedAt, 0),
	}
	if g.ProcessedAt != nil {
		processedAt := time.Unix(*g.ProcessedAt, 0)
		adjustment.ProcessedAt = &processedAt
	}
	return adjustment
}

func (r *GormStockAdjustmentRepository) mapToDomainSlice(gAdjustments []gormStockAdjustmentPending) []domain.StockAdjustmentPending {
	out := make([]domain.StockAdjustmentPending, len(gAdjustments))
	for i, g := range gAdjustments {
		out[i] = *r.mapToDomain(&g)
	}
	return out
}
