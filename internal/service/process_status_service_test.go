package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rixtrayker/getemps-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByNationalNumber(ctx context.Context, nationalNumber string) (*models.User, error) {
	args := m.Called(ctx, nationalNumber)
	return args.Get(0).(*models.User), args.Error(1)
}

// MockSalaryRepository is a mock implementation of SalaryRepository
type MockSalaryRepository struct {
	mock.Mock
}

func (m *MockSalaryRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Salary, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Salary), args.Error(1)
}

func (m *MockSalaryRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

// MockCache is a mock implementation of Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (*models.EmployeeInfo, bool) {
	args := m.Called(key)
	return args.Get(0).(*models.EmployeeInfo), args.Bool(1)
}

func (m *MockCache) Set(key string, value *models.EmployeeInfo, duration time.Duration) {
	m.Called(key, value, duration)
}

func (m *MockCache) Delete(key string) {
	m.Called(key)
}

func TestProcessStatusService_GetEmployeeStatus_WithCache(t *testing.T) {
	// Setup mocks
	mockUserRepo := new(MockUserRepository)
	mockSalaryRepo := new(MockSalaryRepository)
	mockCache := new(MockCache)

	// Create service
	service := NewProcessStatusService(
		mockUserRepo,
		mockSalaryRepo,
		mockCache,
		5*time.Minute,
	)

	ctx := context.Background()
	nationalNumber := "NAT1001"

	t.Run("Cache Hit - Returns Cached Data", func(t *testing.T) {
		// Setup test data
		cachedEmployee := &models.EmployeeInfo{
			ID:             1,
			Username:       "cached_user",
			NationalNumber: nationalNumber,
			Email:          "cached@example.com",
			Phone:          "123456789",
			IsActive:       true,
			SalaryDetails: models.SalaryDetails{
				AverageSalary: 1500.0,
				HighestSalary: 2000.0,
				SumOfSalaries: 9000.0,
			},
			Status:      "RED",
			LastUpdated: time.Now(),
		}

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1001").Return(cachedEmployee, true)

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, cachedEmployee.Username, result.Username)
		assert.Equal(t, cachedEmployee.NationalNumber, result.NationalNumber)
		assert.Equal(t, cachedEmployee.Status, result.Status)

		// Verify mock expectations
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "GetByNationalNumber")
		mockSalaryRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("Cache Miss - Fetches from Database and Caches Result", func(t *testing.T) {
		// Reset mocks
		mockUserRepo.ExpectedCalls = nil
		mockSalaryRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil

		// Setup test data
		user := &models.User{
			ID:             1,
			Username:       "test_user",
			NationalNumber: nationalNumber,
			Email:          "test@example.com",
			Phone:          "123456789",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		salaries := []models.Salary{
			{ID: 1, Year: 2024, Month: 1, Salary: 1000.0, UserID: 1, CreatedAt: time.Now()},
			{ID: 2, Year: 2024, Month: 2, Salary: 1500.0, UserID: 1, CreatedAt: time.Now()},
			{ID: 3, Year: 2024, Month: 3, Salary: 2000.0, UserID: 1, CreatedAt: time.Now()},
		}

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1001").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(user, nil)
		mockSalaryRepo.On("GetByUserID", ctx, user.ID).Return(salaries, nil)
		mockSalaryRepo.On("CountByUserID", ctx, user.ID).Return(len(salaries), nil)
		mockCache.On("Set", "emp_status:NAT1001", mock.AnythingOfType("*models.EmployeeInfo"), 5*time.Minute).Return()

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.NationalNumber, result.NationalNumber)
		assert.Equal(t, "RED", result.Status) // Based on calculated average < 2000

		// Verify salary calculations
		assert.Equal(t, 1500.0, result.SalaryDetails.AverageSalary)
		assert.Equal(t, 2000.0, result.SalaryDetails.HighestSalary)
		assert.Equal(t, 4500.0, result.SalaryDetails.SumOfSalaries)

		// Verify all mocks were called
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertExpectations(t)
	})

	t.Run("Cache Miss - User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo.ExpectedCalls = nil
		mockSalaryRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1001").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return((*models.User)(nil), errors.New("user not found"))

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*AppError)
		assert.True(t, ok)
		assert.Equal(t, 404, appErr.Code)
		assert.Equal(t, "Invalid National Number", appErr.Message)

		// Verify mocks
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("Cache Miss - User Not Active", func(t *testing.T) {
		// Reset mocks
		mockUserRepo.ExpectedCalls = nil
		mockSalaryRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil

		// Setup inactive user
		inactiveUser := &models.User{
			ID:             1,
			Username:       "inactive_user",
			NationalNumber: nationalNumber,
			Email:          "inactive@example.com",
			Phone:          "123456789",
			IsActive:       false, // Inactive user
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1001").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(inactiveUser, nil)

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*AppError)
		assert.True(t, ok)
		assert.Equal(t, 406, appErr.Code) // Not Acceptable

		// Verify mocks
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("Cache Miss - No Salary Data", func(t *testing.T) {
		// Reset mocks
		mockUserRepo.ExpectedCalls = nil
		mockSalaryRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil

		// Setup user
		user := &models.User{
			ID:             1,
			Username:       "no_salary_user",
			NationalNumber: nationalNumber,
			Email:          "nosalary@example.com",
			Phone:          "123456789",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1001").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(user, nil)
		mockSalaryRepo.On("GetByUserID", ctx, user.ID).Return([]models.Salary{}, nil)
		mockSalaryRepo.On("CountByUserID", ctx, user.ID).Return(0, nil)

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		appErr, ok := err.(*AppError)
		assert.True(t, ok)
		assert.Equal(t, 422, appErr.Code) // Unprocessable Entity

		// Verify mocks
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertExpectations(t)
	})
}

func TestProcessStatusService_GetEmployeeStatus_WithoutCache(t *testing.T) {
	// Setup mocks (no cache)
	mockUserRepo := new(MockUserRepository)
	mockSalaryRepo := new(MockSalaryRepository)

	// Create service without cache
	service := NewProcessStatusService(
		mockUserRepo,
		mockSalaryRepo,
		nil, // No cache
		5*time.Minute,
	)

	ctx := context.Background()
	nationalNumber := "NAT1002"

	t.Run("No Cache - Always Fetches from Database", func(t *testing.T) {
		// Setup test data
		user := &models.User{
			ID:             2,
			Username:       "no_cache_user",
			NationalNumber: nationalNumber,
			Email:          "nocache@example.com",
			Phone:          "987654321",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		salaries := []models.Salary{
			{ID: 1, Year: 2024, Month: 1, Salary: 2500.0, UserID: 2, CreatedAt: time.Now()},
			{ID: 2, Year: 2024, Month: 2, Salary: 2500.0, UserID: 2, CreatedAt: time.Now()},
			{ID: 3, Year: 2024, Month: 3, Salary: 2500.0, UserID: 2, CreatedAt: time.Now()},
		}

		// Mock expectations
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(user, nil)
		mockSalaryRepo.On("GetByUserID", ctx, user.ID).Return(salaries, nil)
		mockSalaryRepo.On("CountByUserID", ctx, user.ID).Return(len(salaries), nil)

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.NationalNumber, result.NationalNumber)
		assert.Equal(t, "GREEN", result.Status) // Average = 2500 > 2000

		// Verify salary calculations
		assert.Equal(t, 2500.0, result.SalaryDetails.AverageSalary)
		assert.Equal(t, 2500.0, result.SalaryDetails.HighestSalary)
		assert.Equal(t, 7500.0, result.SalaryDetails.SumOfSalaries)

		// Verify mocks
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertExpectations(t)
	})
}

func TestProcessStatusService_CacheKeyGeneration(t *testing.T) {
	// Setup mocks
	mockUserRepo := new(MockUserRepository)
	mockSalaryRepo := new(MockSalaryRepository)
	mockCache := new(MockCache)

	service := NewProcessStatusService(
		mockUserRepo,
		mockSalaryRepo,
		mockCache,
		5*time.Minute,
	)

	ctx := context.Background()

	testCases := []struct {
		name           string
		nationalNumber string
		expectedKey    string
	}{
		{
			name:           "Normal National Number",
			nationalNumber: "NAT1001",
			expectedKey:    "emp_status:NAT1001",
		},
		{
			name:           "Empty National Number",
			nationalNumber: "",
			expectedKey:    "emp_status:",
		},
		{
			name:           "Special Characters",
			nationalNumber: "NAT@123#$%",
			expectedKey:    "emp_status:NAT@123#$%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks
			mockCache.ExpectedCalls = nil

			// Mock cache miss
			mockCache.On("Get", tc.expectedKey).Return((*models.EmployeeInfo)(nil), false)
			mockUserRepo.On("GetByNationalNumber", ctx, tc.nationalNumber).Return((*models.User)(nil), errors.New("user not found"))

			// Execute (will fail, but we just want to verify the cache key)
			_, _ = service.GetEmployeeStatus(ctx, tc.nationalNumber)

			// Verify the correct cache key was used
			mockCache.AssertCalled(t, "Get", tc.expectedKey)
		})
	}
}

func TestProcessStatusService_CacheIntegration_BusinessLogic(t *testing.T) {
	// Setup mocks
	mockUserRepo := new(MockUserRepository)
	mockSalaryRepo := new(MockSalaryRepository)
	mockCache := new(MockCache)

	service := NewProcessStatusService(
		mockUserRepo,
		mockSalaryRepo,
		mockCache,
		5*time.Minute,
	)

	ctx := context.Background()

	t.Run("Complex Salary Calculation with Caching", func(t *testing.T) {
		nationalNumber := "NAT1005"
		
		// Setup user
		user := &models.User{
			ID:             5,
			Username:       "complex_user",
			NationalNumber: nationalNumber,
			Email:          "complex@example.com",
			Phone:          "123456789",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Setup salaries with seasonal adjustments and tax implications
		salaries := []models.Salary{
			{ID: 1, Year: 2024, Month: 6, Salary: 5000.0, UserID: 5, CreatedAt: time.Now()},  // Summer month (-5%)
			{ID: 2, Year: 2024, Month: 12, Salary: 5000.0, UserID: 5, CreatedAt: time.Now()}, // December (+10%)
			{ID: 3, Year: 2024, Month: 3, Salary: 5000.0, UserID: 5, CreatedAt: time.Now()},  // Regular month
		}

		// Mock expectations
		mockCache.On("Get", "emp_status:NAT1005").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(user, nil)
		mockSalaryRepo.On("GetByUserID", ctx, user.ID).Return(salaries, nil)
		mockSalaryRepo.On("CountByUserID", ctx, user.ID).Return(len(salaries), nil)
		mockCache.On("Set", "emp_status:NAT1005", mock.AnythingOfType("*models.EmployeeInfo"), 5*time.Minute).Return()

		// Execute
		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		// Assertions
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify business logic calculations
		// Expected: 4750 + 5500 + 5000 = 15250
		// Average: 15250 / 3 = 5083.33
		// Tax deduction (7% because > 10000): 15250 * 0.93 = 14182.5
		// Final average: 14182.5 / 3 = 4727.5
		
		expectedSum := 14182.5 // After tax deduction
		expectedAvg := expectedSum / 3
		
		assert.Equal(t, "GREEN", result.Status) // > 2000
		assert.InDelta(t, expectedAvg, result.SalaryDetails.AverageSalary, 0.1)
		assert.InDelta(t, expectedSum, result.SalaryDetails.SumOfSalaries, 0.1)

		// Verify mocks
		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertExpectations(t)
	})
}

func TestProcessStatusService_ErrorHandling(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockSalaryRepo := new(MockSalaryRepository)
	mockCache := new(MockCache)

	service := NewProcessStatusService(
		mockUserRepo,
		mockSalaryRepo,
		mockCache,
		5*time.Minute,
	)

	ctx := context.Background()
	nationalNumber := "NAT_ERROR"

	t.Run("Database Error During User Fetch", func(t *testing.T) {
		mockCache.On("Get", "emp_status:NAT_ERROR").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return((*models.User)(nil), errors.New("database connection failed"))

		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to fetch user")

		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Database Error During Salary Fetch", func(t *testing.T) {
		// Reset mocks
		mockUserRepo.ExpectedCalls = nil
		mockSalaryRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil

		user := &models.User{
			ID:             1,
			Username:       "error_user",
			NationalNumber: nationalNumber,
			Email:          "error@example.com",
			Phone:          "123456789",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockCache.On("Get", "emp_status:NAT_ERROR").Return((*models.EmployeeInfo)(nil), false)
		mockUserRepo.On("GetByNationalNumber", ctx, nationalNumber).Return(user, nil)
		mockSalaryRepo.On("CountByUserID", ctx, user.ID).Return(0, errors.New("salary table corrupted"))

		result, err := service.GetEmployeeStatus(ctx, nationalNumber)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to count salary records")

		mockCache.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockSalaryRepo.AssertExpectations(t)
	})
}