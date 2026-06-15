package domain

// ProductIngredient define a ficha técnica:
// qual ingrediente e em qual quantidade compõe um produto.
// Exemplo: Pizza Margherita usa 0.150 kg de Queijo Mozzarella.
type ProductIngredient struct {
	ID           uint
	ProductID    uint
	IngredientID uint
	Quantity     float64 // quantidade consumida por unidade do produto

	// Preenchido em joins (leitura)
	Ingredient *Ingredient
}
