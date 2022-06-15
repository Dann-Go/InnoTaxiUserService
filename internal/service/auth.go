package service

import (
	"errors"
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	signingKey = "grk#iwoerjn%324lskdfnHsdj3skldf"
	tokenTTL   = 48 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

type AuthService struct {
	repo repository.Authorization
}

func (s *AuthService) CreateUser(user *domain.User) (*domain.User, error) {
	user.PasswordHash = hashPassword(user.PasswordHash)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(phone, password string) (string, error) {
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return "", err
	}
	if isValid := checkPasswordHash(password, user.PasswordHash); !isValid {
		err := errors.New("Wrong Password")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) GetUserByPhone(phone string) (*domain.User, error) {
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, err
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
