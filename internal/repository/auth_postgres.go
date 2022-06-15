package repository

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users(name, phone, email, password_hash) VALUES ($1, $2, $3, $4);`
	_, err := r.db.Exec(query, user.Name, user.Phone, user.Email, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	user, err = r.GetUserByPhone(user.Phone)
	user.PasswordHash = ""
	return user, err
}

func (r *AuthPostgres) GetUserByPhone(phone string) (*domain.User, error) {
	user := domain.User{}
	query := `SELECT * FROM users WHERE phone=$1 `
	err := r.db.Get(&user, query, phone)
	return &user, err
}
