package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rixtrayker/getemps-service/internal/models"
)

type salaryRepository struct {
	db *sqlx.DB
}

func NewSalaryRepository(db *sqlx.DB) SalaryRepository {
	return &salaryRepository{db: db}
}

func (r *salaryRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Salary, error) {
	query := `
		SELECT id, year, month, salary, user_id, created_at
		FROM salaries
		WHERE user_id = $1
		ORDER BY year ASC, month ASC
	`

	var salaries []models.Salary
	err := r.db.SelectContext(ctx, &salaries, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get salaries for user %d: %w", userID, err)
	}

	return salaries, nil
}

func (r *salaryRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM salaries
		WHERE user_id = $1
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count salaries for user %d: %w", userID, err)
	}

	return count, nil
}
