package repository

import (
	"context"

	"github.com/yourusername/getemps-service/internal/models"
)

type UserRepository interface {
	GetByNationalNumber(ctx context.Context, nationalNumber string) (*models.User, error)
}

type SalaryRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]models.Salary, error)
	CountByUserID(ctx context.Context, userID int64) (int, error)
}