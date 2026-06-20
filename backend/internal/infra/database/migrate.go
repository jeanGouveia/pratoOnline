package database

import (
	"fmt"

	"gorm.io/gorm"
)

// gormMigrationModel espelha a tabela criada pelo Goose para não colidir.
// Usamos AutoMigrate apenas nas tabelas de domínio — o schema real fica nas migrations SQL.
// Este arquivo garante que as tabelas existam no SQLite em dev sem precisar rodar Goose manualmente.

type userMigrationModel struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
}

func (userMigrationModel) TableName() string { return "users" }

type productMigrationModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null;default:0"`
	IsComposto  bool    `gorm:"not null;default:false"`
	Active      bool    `gorm:"not null;default:true"`
	CreatedAt   int64   `gorm:"autoCreateTime"`
	UpdatedAt   int64   `gorm:"autoUpdateTime"`
}

func (productMigrationModel) TableName() string { return "products" }

type ingredientMigrationModel struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"not null"`
	Unit          string  `gorm:"not null"` // ex: "kg", "L", "un"
	StockQuantity float64 `gorm:"not null;default:0"`
	MinStock      float64 `gorm:"not null;default:0"`
	CreatedAt     int64   `gorm:"autoCreateTime"`
	UpdatedAt     int64   `gorm:"autoUpdateTime"`
}

func (ingredientMigrationModel) TableName() string { return "ingredients" }

type productIngredientMigrationModel struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	ProductID    uint    `gorm:"not null;index"`
	IngredientID uint    `gorm:"not null"`
	Quantity     float64 `gorm:"not null"`
	CreatedAt    int64   `gorm:"autoCreateTime"`
	UpdatedAt    int64   `gorm:"autoUpdateTime"`
}

func (productIngredientMigrationModel) TableName() string { return "product_ingredients" }

type productCompositionMigrationModel struct {
	ID                 uint    `gorm:"primaryKey;autoIncrement"`
	ParentProductID    uint    `gorm:"not null;index"`
	ComponentProductID uint    `gorm:"not null"`
	Quantity           float64 `gorm:"not null"`
	CreatedAt          int64   `gorm:"autoCreateTime"`
	UpdatedAt          int64   `gorm:"autoUpdateTime"`
}

func (productCompositionMigrationModel) TableName() string { return "product_compositions" }

type orderMigrationModel struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	Status     string  `gorm:"not null;default:'pending'"`
	TotalPrice float64 `gorm:"not null;default:0"`
	Notes      string
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
}

func (orderMigrationModel) TableName() string { return "orders" }

type orderItemMigrationModel struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null"`
	Quantity  float64 `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
	CreatedAt int64   `gorm:"autoCreateTime"`
	UpdatedAt int64   `gorm:"autoUpdateTime"`
}

func (orderItemMigrationModel) TableName() string { return "order_items" }

type stockAdjustmentPendingMigrationModel struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	OrderID      uint    `gorm:"not null;index"`
	IngredientID uint    `gorm:"not null;index"`
	Quantity     float64 `gorm:"not null"`
	OrderStatus  string  `gorm:"not null"`
	Status       string  `gorm:"not null;default:'pending';index"`
	CreatedAt    int64   `gorm:"autoCreateTime"`
	ProcessedAt  *int64  `gorm:"index"`
}

func (stockAdjustmentPendingMigrationModel) TableName() string { return "stock_adjustments_pending" }

// RunMigrations executa o AutoMigrate do GORM para todas as tabelas do sistema.
// Em produção com Oracle, substitua por Goose com migrations SQL versionadas.
func RunMigrations(db *gorm.DB) error {
	models := []interface{}{
		&userMigrationModel{},
		&productMigrationModel{},
		&ingredientMigrationModel{},
		&productIngredientMigrationModel{},
		&productCompositionMigrationModel{},
		&orderMigrationModel{},
		&orderItemMigrationModel{},
		&stockAdjustmentPendingMigrationModel{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}

	return nil
}
