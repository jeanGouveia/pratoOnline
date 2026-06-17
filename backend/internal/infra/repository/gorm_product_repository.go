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

type gormProduct struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null;default:0"`
	IsComposto  bool    `gorm:"not null;default:false"`
	Active      bool    `gorm:"not null;default:true"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

func (gormProduct) TableName() string { return "products" }

type gormIngredient struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"not null"`
	Unit          string  `gorm:"not null;default:'un'"`
	StockQuantity float64 `gorm:"not null;default:0"`
	MinStock      float64 `gorm:"not null;default:0"`
	CreatedAt     int64   `gorm:"autoCreateTime"`
	UpdatedAt     int64   `gorm:"autoUpdateTime"`
}

func (gormIngredient) TableName() string { return "ingredients" }

type gormProductIngredient struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	ProductID    uint           `gorm:"not null;index"`
	IngredientID uint           `gorm:"not null"`
	Quantity     float64        `gorm:"not null"`
	Ingredient   gormIngredient `gorm:"foreignKey:IngredientID"`
}

func (gormProductIngredient) TableName() string { return "product_ingredients" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.ProductRepository = (*GormProductRepository)(nil)

type GormProductRepository struct{ db *gorm.DB }

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	m := gormProduct{
		Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateProduct: %w", err)
	}
	p.ID = m.ID
	p.CreatedAt = time.Unix(m.CreatedAt, 0)
	p.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	var m gormProduct
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindProductByID: %w", err)
	}
	return productToDomain(&m), nil
}

func (r *GormProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var ms []gormProduct
	if err := r.db.WithContext(ctx).Where("active = ?", true).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListProducts: %w", err)
	}
	out := make([]domain.Product, len(ms))
	for i, m := range ms {
		out[i] = *productToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	m := gormProduct{
		ID: p.ID, Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateProduct: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	// Soft delete: marca Active=false
	if err := r.db.WithContext(ctx).Model(&gormProduct{}).
		Where("id = ?", id).Update("active", false).Error; err != nil {
		return fmt.Errorf("DeleteProduct: %w", err)
	}
	return nil
}

// ── Ingrediente ──────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateIngredient(ctx context.Context, i *domain.Ingredient) error {
	m := gormIngredient{
		Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateIngredient: %w", err)
	}
	i.ID = m.ID
	i.CreatedAt = time.Unix(m.CreatedAt, 0)
	i.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindIngredientByID(ctx context.Context, id uint) (*domain.Ingredient, error) {
	var m gormIngredient
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindIngredientByID: %w", err)
	}
	return ingredientToDomain(&m), nil
}

func (r *GormProductRepository) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ms []gormIngredient
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListIngredients: %w", err)
	}
	out := make([]domain.Ingredient, len(ms))
	for i, m := range ms {
		out[i] = *ingredientToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient) error {
	m := gormIngredient{
		ID: i.ID, Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateIngredient: %w", err)
	}
	return nil
}

// ── Ficha técnica ────────────────────────────────────────────────────────────

func (r *GormProductRepository) SetProductIngredients(
	ctx context.Context, productID uint, items []domain.ProductIngredient,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Apaga ficha anterior e recria (upsert simples)
		if err := tx.Where("product_id = ?", productID).
			Delete(&gormProductIngredient{}).Error; err != nil {
			return fmt.Errorf("SetProductIngredients delete: %w", err)
		}
		for _, item := range items {
			m := gormProductIngredient{
				ProductID:    productID,
				IngredientID: item.IngredientID,
				Quantity:     item.Quantity,
			}
			if err := tx.Create(&m).Error; err != nil {
				return fmt.Errorf("SetProductIngredients insert: %w", err)
			}
		}
		return nil
	})
}

func (r *GormProductRepository) GetProductIngredients(
	ctx context.Context, productID uint,
) ([]domain.ProductIngredient, error) {
	var ms []gormProductIngredient
	if err := r.db.WithContext(ctx).
		Preload("Ingredient").
		Where("product_id = ?", productID).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("GetProductIngredients: %w", err)
	}
	out := make([]domain.ProductIngredient, len(ms))
	for i, m := range ms {
		ing := ingredientToDomain(&m.Ingredient)
		out[i] = domain.ProductIngredient{
			ID: m.ID, ProductID: m.ProductID,
			IngredientID: m.IngredientID, Quantity: m.Quantity,
			Ingredient: ing,
		}
	}
	return out, nil
}

// ── Estoque ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) DecreaseIngredientStock(
	ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB,
) error {
	// Usa o DB da transação se fornecido, senão usa o DB padrão
	db := r.db
	if txDB != nil {
		db = txDB.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	// Usa UPDATE com CHECK inline para garantir que não vai negativo
	result := db.
		Model(&gormIngredient{}).
		Where("id = ? AND stock_quantity >= ?", ingredientID, qty).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))

	if result.Error != nil {
		return fmt.Errorf("DecreaseIngredientStock id=%d: %w", ingredientID, result.Error)
	}
	if result.RowsAffected == 0 {
		// Busca o nome para mensagem de erro amigável
		var ing gormIngredient
		log.Printf("DecreaseIngredientStock - ingredientID recebido: %d", ingredientID)
		err := r.db.WithContext(ctx).Select("name, stock_quantity").First(&ing, ingredientID).Error
		log.Printf("DecreaseIngredientStock - erro da query: %v", err)
		log.Printf("DecreaseIngredientStock - ing.ID: %d", ing.ID)
		log.Printf("DecreaseIngredientStock - ing.Name: '%s'", ing.Name)
		log.Printf("DecreaseIngredientStock - ing.StockQuantity: %f", ing.StockQuantity)
		return fmt.Errorf(
			"estoque insuficiente para '%s': disponível=%.4f necessário=%.4f",
			ing.Name, ing.StockQuantity, qty,
		)
	}
	return nil
}

// ── Mappers ──────────────────────────────────────────────────────────────────

func productToDomain(m *gormProduct) *domain.Product {
	return &domain.Product{
		ID: m.ID, Name: m.Name, Description: m.Description,
		Price: m.Price, IsComposto: m.IsComposto, Active: m.Active,
		CreatedAt: time.Unix(m.CreatedAt, 0), UpdatedAt: time.Unix(m.UpdatedAt, 0),
	}
}

func ingredientToDomain(m *gormIngredient) *domain.Ingredient {
	return &domain.Ingredient{
		ID: m.ID, Name: m.Name, Unit: m.Unit,
		StockQuantity: m.StockQuantity, MinStock: m.MinStock,
		CreatedAt: time.Unix(m.CreatedAt, 0), UpdatedAt: time.Unix(m.UpdatedAt, 0),
	}
}
