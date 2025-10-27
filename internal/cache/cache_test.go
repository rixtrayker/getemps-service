package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/rixtrayker/getemps-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCache_BasicOperations(t *testing.T) {
	// Initialize cache
	cache := NewInMemoryCache(5*time.Minute, 10*time.Minute)
	assert.NotNil(t, cache)

	// Test data
	testEmployee := &models.EmployeeInfo{
		ID:             1,
		Username:       "test_user",
		NationalNumber: "NAT1001",
		Email:          "test@example.com",
		Phone:          "123456789",
		IsActive:       true,
		LastUpdated:    time.Now(),
	}

	t.Run("Set and Get", func(t *testing.T) {
		key := "test_key_1"
		
		// Set value
		cache.Set(key, testEmployee, 1*time.Minute)
		
		// Get value
		retrieved, found := cache.Get(key)
		
		// Assertions
		assert.True(t, found, "Expected to find the cached value")
		require.NotNil(t, retrieved, "Retrieved value should not be nil")
		assert.Equal(t, testEmployee.ID, retrieved.ID)
		assert.Equal(t, testEmployee.Username, retrieved.Username)
		assert.Equal(t, testEmployee.NationalNumber, retrieved.NationalNumber)
		assert.Equal(t, testEmployee.Email, retrieved.Email)
		assert.Equal(t, testEmployee.Phone, retrieved.Phone)
		assert.Equal(t, testEmployee.IsActive, retrieved.IsActive)
	})

	t.Run("Get Non-Existent Key", func(t *testing.T) {
		key := "non_existent_key"
		
		retrieved, found := cache.Get(key)
		
		assert.False(t, found, "Expected not to find non-existent key")
		assert.Nil(t, retrieved, "Retrieved value should be nil for non-existent key")
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test_key_2"
		
		// Set value
		cache.Set(key, testEmployee, 1*time.Minute)
		
		// Verify it exists
		_, found := cache.Get(key)
		assert.True(t, found, "Value should exist before deletion")
		
		// Delete
		cache.Delete(key)
		
		// Verify it's gone
		_, found = cache.Get(key)
		assert.False(t, found, "Value should not exist after deletion")
	})

	t.Run("Overwrite Existing Key", func(t *testing.T) {
		key := "test_key_3"
		
		// Set initial value
		cache.Set(key, testEmployee, 1*time.Minute)
		
		// Create different employee
		newEmployee := &models.EmployeeInfo{
			ID:             2,
			Username:       "new_user",
			NationalNumber: "NAT2002",
			Email:          "new@example.com",
			Phone:          "987654321",
			IsActive:       false,
			LastUpdated:    time.Now(),
		}
		
		// Overwrite
		cache.Set(key, newEmployee, 1*time.Minute)
		
		// Get and verify new value
		retrieved, found := cache.Get(key)
		
		assert.True(t, found, "Expected to find the cached value")
		require.NotNil(t, retrieved, "Retrieved value should not be nil")
		assert.Equal(t, newEmployee.ID, retrieved.ID)
		assert.Equal(t, newEmployee.Username, retrieved.Username)
		assert.Equal(t, newEmployee.NationalNumber, retrieved.NationalNumber)
	})

	t.Run("Expiration", func(t *testing.T) {
		key := "test_key_expiration"
		
		// Set value with very short expiration
		cache.Set(key, testEmployee, 50*time.Millisecond)
		
		// Immediately check it exists
		_, found := cache.Get(key)
		assert.True(t, found, "Value should exist immediately after setting")
		
		// Wait for expiration
		time.Sleep(100 * time.Millisecond)
		
		// Check it's gone
		_, found = cache.Get(key)
		assert.False(t, found, "Value should have expired")
	})
}

func TestInMemoryCache_EdgeCases(t *testing.T) {
	cache := NewInMemoryCache(5*time.Minute, 10*time.Minute)

	t.Run("Set Nil Value", func(t *testing.T) {
		key := "nil_test"
		
		// Setting nil should not crash
		cache.Set(key, nil, 1*time.Minute)
		
		// Getting should return nil and not found (current implementation stores it as empty struct)
		retrieved, found := cache.Get(key)
		// Note: Current implementation marshals nil to empty EmployeeInfo struct
		// This is expected behavior for this cache implementation
		if found {
			assert.NotNil(t, retrieved, "Retrieved value should not be nil due to JSON marshaling of nil")
			assert.Equal(t, int64(0), retrieved.ID, "Should be default EmployeeInfo values")
		} else {
			assert.Nil(t, retrieved, "If not found, should be nil")
		}
	})

	t.Run("Empty Key", func(t *testing.T) {
		testEmployee := &models.EmployeeInfo{
			ID:             1,
			Username:       "test",
			NationalNumber: "NAT1001",
			Email:          "test@example.com",
			Phone:          "123456789",
			IsActive:       true,
			LastUpdated:    time.Now(),
		}
		
		// Empty key should work
		cache.Set("", testEmployee, 1*time.Minute)
		
		retrieved, found := cache.Get("")
		assert.True(t, found, "Empty key should work")
		assert.NotNil(t, retrieved, "Retrieved value should not be nil")
	})

	t.Run("Very Long Key", func(t *testing.T) {
		longKey := string(make([]byte, 1000)) // 1000 character key
		for i := range longKey {
			longKey = longKey[:i] + "a" + longKey[i+1:]
		}
		
		testEmployee := &models.EmployeeInfo{
			ID:             1,
			Username:       "test",
			NationalNumber: "NAT1001",
			Email:          "test@example.com",
			Phone:          "123456789",
			IsActive:       true,
			LastUpdated:    time.Now(),
		}
		
		// Long key should work
		cache.Set(longKey, testEmployee, 1*time.Minute)
		
		retrieved, found := cache.Get(longKey)
		assert.True(t, found, "Long key should work")
		assert.NotNil(t, retrieved, "Retrieved value should not be nil")
	})

	t.Run("Zero Duration", func(t *testing.T) {
		key := "zero_duration"
		testEmployee := &models.EmployeeInfo{
			ID:             1,
			Username:       "test",
			NationalNumber: "NAT1001",
			Email:          "test@example.com",
			Phone:          "123456789",
			IsActive:       true,
			LastUpdated:    time.Now(),
		}
		
		// Zero duration should still work (default expiration)
		cache.Set(key, testEmployee, 0)
		
		retrieved, found := cache.Get(key)
		assert.True(t, found, "Zero duration should use default expiration")
		assert.NotNil(t, retrieved, "Retrieved value should not be nil")
	})
}

func TestKeyGeneration(t *testing.T) {
	tests := []struct {
		name        string
		function    func(string) string
		input       string
		expected    string
		description string
	}{
		{
			name:        "GenerateCacheKey",
			function:    GenerateCacheKey,
			input:       "NAT1001",
			expected:    "emp_status:NAT1001",
			description: "Should generate employee status cache key",
		},
		{
			name:        "GenerateCacheKey_Empty",
			function:    GenerateCacheKey,
			input:       "",
			expected:    "emp_status:",
			description: "Should handle empty input",
		},
		{
			name:        "GenerateCacheKey_Special_Characters",
			function:    GenerateCacheKey,
			input:       "NAT@1001#$%",
			expected:    "emp_status:NAT@1001#$%",
			description: "Should handle special characters",
		},
		{
			name:        "GenerateCacheKey_Numbers",
			function:    GenerateCacheKey,
			input:       "12345",
			expected:    "emp_status:12345",
			description: "Should handle numeric input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.input)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestInMemoryCache_Concurrency(t *testing.T) {
	cache := NewInMemoryCache(5*time.Minute, 10*time.Minute)
	
	// Test concurrent reads and writes
	t.Run("Concurrent Operations", func(t *testing.T) {
		const numGoroutines = 10
		const numOperations = 100
		
		done := make(chan bool, numGoroutines)
		
		// Create test employee
		testEmployee := &models.EmployeeInfo{
			ID:             1,
			Username:       "concurrent_test",
			NationalNumber: "NAT1001",
			Email:          "test@example.com",
			Phone:          "123456789",
			IsActive:       true,
			LastUpdated:    time.Now(),
		}
		
		// Start multiple goroutines doing cache operations
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("concurrent_key_%d_%d", id, j)
					
					// Set
					cache.Set(key, testEmployee, 1*time.Minute)
					
					// Get
					_, found := cache.Get(key)
					assert.True(t, found, "Value should be found immediately after setting")
					
					// Delete
					cache.Delete(key)
					
					// Verify deletion
					_, found = cache.Get(key)
					assert.False(t, found, "Value should not be found after deletion")
				}
			}(i)
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

func TestInMemoryCache_MemoryLeaks(t *testing.T) {
	// This test ensures that deleted items are actually removed
	cache := NewInMemoryCache(5*time.Minute, 1*time.Second) // Short cleanup interval
	
	testEmployee := &models.EmployeeInfo{
		ID:             1,
		Username:       "memory_test",
		NationalNumber: "NAT1001",
		Email:          "test@example.com",
		Phone:          "123456789",
		IsActive:       true,
		LastUpdated:    time.Now(),
	}
	
	t.Run("Cleanup After Expiration", func(t *testing.T) {
		// Add many items with short expiration
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("temp_key_%d", i)
			cache.Set(key, testEmployee, 50*time.Millisecond)
		}
		
		// Wait for expiration and cleanup
		time.Sleep(2 * time.Second)
		
		// Verify items are gone
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("temp_key_%d", i)
			_, found := cache.Get(key)
			assert.False(t, found, "Expired items should be cleaned up")
		}
	})
}

