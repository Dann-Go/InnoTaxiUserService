package service

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain/apperrors"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
)

type User interface {
	CreateUser(user *domain.User) (*domain.UserResponse, error)
	GetUserByPhone(phone string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type UserService struct {
	repo *repository.UserRepository
}

func (s *UserService) CreateUser(user *domain.User) (*domain.UserResponse, error) {
	user.PasswordHash = HashPassword(user.PasswordHash)

	if userCheck, err := s.repo.GetUserByPhone(user.Phone); userCheck.Phone != "" {
		return nil, apperrors.Wrapper(apperrors.ErrPhoneIsAlreadyTaken, err)

	} else if userCheck, err = s.repo.GetUserByEmail(user.Email); userCheck.Email != "" {
		return nil, apperrors.Wrapper(apperrors.ErrEmailIsAlreadyTaken, err)
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
