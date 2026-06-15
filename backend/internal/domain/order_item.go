package domain

// OrderItem representa um produto dentro de um pedido.
// UnitPrice é salvo no momento da venda (snapshot — preço pode mudar depois).
type OrderItem struct {
	ID        uint
	OrderID   uint
	ProductID uint
	Quantity  float64
	UnitPrice float64 // snapshot do preço no momento do pedido

	// Preenchido em joins (leitura)
	Product *Product
}
