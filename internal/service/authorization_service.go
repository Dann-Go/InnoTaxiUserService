package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/Dann-Go/InnoTaxiUserService/internal/config"
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain/apperrors"

	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Authorization interface {
	GenerateToken(phone, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type AuthorizationService struct {
	repo *repository.UserRepository
}

func NewAuthorizationService(repo *repository.UserRepository) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

var authCfg = config.NewAuthConfig()

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

func (ms *AuthorizationService) GenerateToken(phone, password string) (string, error) {
	user, err := ms.repo.GetUserByPhone(phone)
	if err != nil {
		return "", err
	}
	if isValid := checkPasswordHash(password, user.PasswordHash); !isValid {
		err := apperrors.Wrapper(apperrors.ErrWrongPassword, err)
		return "", err
	}
	tokenTTL, _ := strconv.Atoi(authCfg.TokenTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(tokenTTL) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(authCfg.SigningKey))
}

func (ms *AuthorizationService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.Wrapper(apperrors.ErrInvalidSigningMethod, errors.New("invalid signing method"))
		}
		return []byte(authCfg.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, apperrors.Wrapper(apperrors.ErrWrongTokenClaims, errors.New("wrong token claims"))
	}

	return claims.UserId, err
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
