package service

import (
	"errors"

	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
)

type User interface {
	CreateUser(user *domain.User) (*domain.UserResponse, error)
	GenerateToken(phone, password string) (string, error)
	GetUserByPhone(phone string) (*domain.User, error)
	ParseToken(token string) (int, error)
}

type UserService struct {
	repo *repository.UserRepository
}

func (s *UserService) CreateUser(user *domain.User) (*domain.UserResponse, error) {
	user.PasswordHash = HashPassword(user.PasswordHash)

	if userCheck, _ := s.repo.GetUserByPhone(user.Phone); userCheck.Phone != "" {
		return nil, errors.New("user with such phone already exists")

	} else if userCheck, _ := s.repo.GetUserByEmail(user.Email); userCheck.Email != "" {
		return nil, errors.New("user with such email already exists")
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *UserService) GetUserByPhone(phone string) (*domain.User, error) {
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	return user, err
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
