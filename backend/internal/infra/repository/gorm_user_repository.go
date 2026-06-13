package repository

import (
	"errors"

	"github.com/seu-usuario/my-app/backend/internal/domain"
	"github.com/seu-usuario/my-app/backend/internal/ports"
	"gorm.io/gorm"
)

// gormUserModel é o modelo GORM. Tags de banco APENAS aqui, nunca no domain.
type gormUserModel struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}

// GormUserRepository implementa ports.UserRepository via GORM.
// Verificação em tempo de compilação — falha se a interface não for satisfeita.
var _ ports.UserRepository = (*GormUserRepository)(nil)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *domain.User) error {
	model := toModel(user)
	result := r.db.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	user.ID = uint(model.ID)
	return nil
}

func (r *GormUserRepository) FindByEmail(email string) (*domain.User, error) {
	var model gormUserModel
	result := r.db.Where("email = ?", email).First(&model)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(&model), nil
}

func (r *GormUserRepository) FindByID(id uint) (*domain.User, error) {
	var model gormUserModel
	result := r.db.First(&model, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(&model), nil
}

// Mappers: isolam a conversão entre domain e infra
func toModel(u *domain.User) gormUserModel {
	return gormUserModel{
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

func toDomain(m *gormUserModel) *domain.User {
	return &domain.User{
		ID:           uint(m.ID),
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
