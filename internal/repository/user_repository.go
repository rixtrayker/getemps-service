package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rixtrayker/getemps-service/internal/models"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByNationalNumber(ctx context.Context, nationalNumber string) (*models.User, error) {
	query := `
		SELECT id, username, national_number, email, phone, is_active, created_at, updated_at
		FROM users
		WHERE national_number = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, nationalNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with national number %s not found", nationalNumber)
		}
		return nil, fmt.Errorf("failed to get user by national number: %w", err)
	}

	return &user, nil
}