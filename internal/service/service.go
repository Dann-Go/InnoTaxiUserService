package service

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
)

type Authorization interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GenerateToken(phone, password string) (string, error)
	GetUserByPhone(phone string) (*domain.User, error)
	ParseToken(token string) (int, error)
}

type User interface {
	OrderTaxi()
	RateDrive()
	GetRating()
	GetAllDrives()
}

type Service struct {
	Authorization
	User
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization),
	}
}
