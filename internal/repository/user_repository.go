package repository

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/domain"
	"github.com/jmoiron/sqlx"
)

type User interface {
	CreateUser(user *domain.User) (*domain.UserResponse, error)
	GetUserByPhone(phone string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

const createUserQuery = `	INSERT INTO users(name, phone, email, password_hash) 
								VALUES ($1, $2, $3, $4)
							RETURNING
    							id,name, phone, email, rating;`

func (r *UserRepository) CreateUser(user *domain.User) (*domain.UserResponse, error) {

	row := r.db.QueryRow(createUserQuery, user.Name, user.Phone, user.Email, user.PasswordHash)
	var userResponse domain.UserResponse

	err := row.Scan(
		&userResponse.ID,
		&userResponse.Name,
		&userResponse.Email,
		&userResponse.Phone,
		&userResponse.Rating)
	if err != nil {
		return nil, err
	}

	return &userResponse, err
}

func (r *UserRepository) GetUserByPhone(phone string) (*domain.User, error) {
	user := domain.User{}
	query := `SELECT * FROM users WHERE phone=$1 `
	err := r.db.Get(&user, query, phone)
	return &user, err
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	query := `SELECT * FROM users WHERE email=$1 `
	err := r.db.Get(&user, query, email)
	return &user, err
}
