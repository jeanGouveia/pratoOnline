package domain

import "time"

type Ingredient struct {
	ID            uint
	Name          string
	Unit          string  // "kg", "L", "un", "g", "ml"
	StockQuantity float64 // quantidade atual em estoque
	MinStock      float64 // alerta de estoque mínimo (opcional)
	Active        bool    // soft delete flag
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
