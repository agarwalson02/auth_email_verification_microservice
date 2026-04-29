package repository

import (
	"auth_service/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, first_name, last_name, password, role, avatar)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Role,
		user.Avatar,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT uuid, email, first_name, last_name, password, role, avatar, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT uuid, email, first_name, last_name, password, role, avatar, created_at, updated_at
		FROM users
		WHERE uuid = $1
	`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
