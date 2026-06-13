package ports

import "github.com/seu-usuario/my-app/backend/internal/domain"

// UserRepository é o contrato que isola o domínio da infraestrutura.
// Qualquer banco de dados (SQLite, Oracle, Postgres) deve implementar esta interface.
// NENHUMA query SQL, NENHUM tipo GORM deve aparecer fora de internal/infra/repository/.
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}
