package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rixtrayker/getemps-service/internal/cache"
	"github.com/rixtrayker/getemps-service/internal/models"
	"github.com/rixtrayker/getemps-service/internal/repository"
)

type ProcessStatusService struct {
	userRepo   repository.UserRepository
	salaryRepo repository.SalaryRepository
	calculator *SalaryCalculator
	cache      cache.Cache
	cacheTTL   time.Duration
}

func NewProcessStatusService(
	userRepo repository.UserRepository,
	salaryRepo repository.SalaryRepository,
	cache cache.Cache,
	cacheTTL time.Duration,
) *ProcessStatusService {
	return &ProcessStatusService{
		userRepo:   userRepo,
		salaryRepo: salaryRepo,
		calculator: NewSalaryCalculator(),
		cache:      cache,
		cacheTTL:   cacheTTL,
	}
}

func (s *ProcessStatusService) GetEmployeeStatus(ctx context.Context, nationalNumber string) (*models.EmployeeInfo, error) {
	// Step 1: Check cache first
	if s.cache != nil {
		cacheKey := cache.GenerateCacheKey(nationalNumber)
		if cachedResult, found := s.cache.Get(cacheKey); found {
			return cachedResult, nil
		}
	}

	// Step 2: Validate and fetch user
	user, err := s.userRepo.GetByNationalNumber(ctx, nationalNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, &AppError{
				Code:    404,
				Message: "Invalid National Number",
			}
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Step 3: Check if user is active
	if !user.IsActive {
		return nil, &AppError{
			Code:    406,
			Message: "User is not Active",
		}
	}

	// Step 4: Check salary records count
	salaryCount, err := s.salaryRepo.CountByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count salary records: %w", err)
	}

	if salaryCount < 3 {
		return nil, &AppError{
			Code:    422,
			Message: "INSUFFICIENT_DATA",
		}
	}

	// Step 5: Fetch salary records
	salaries, err := s.salaryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch salary records: %w", err)
	}

	// Step 6: Calculate salary statistics and status
	calculation := s.calculator.CalculateEmployeeStatus(salaries)

	// Step 7: Build response
	employeeInfo := &models.EmployeeInfo{
		ID:             user.ID,
		Username:       user.Username,
		NationalNumber: user.NationalNumber,
		Email:          user.Email,
		Phone:          user.Phone,
		IsActive:       user.IsActive,
		SalaryDetails: models.SalaryDetails{
			AverageSalary: calculation.AverageSalary,
			HighestSalary: calculation.HighestSalary,
			SumOfSalaries: calculation.SumOfSalaries,
		},
		Status:      calculation.Status,
		LastUpdated: time.Now(),
	}

	// Step 8: Cache the result
	if s.cache != nil {
		cacheKey := cache.GenerateCacheKey(nationalNumber)
		s.cache.Set(cacheKey, employeeInfo, s.cacheTTL)
	}

	return employeeInfo, nil
}

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}
