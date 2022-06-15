package repository

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByPhone(phone string) (*domain.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
