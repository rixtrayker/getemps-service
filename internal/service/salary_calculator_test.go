package service

import (
	"testing"

	"github.com/yourusername/getemps-service/internal/models"
)

func TestSalaryCalculator_ApplySeasonalAdjustments(t *testing.T) {
	calculator := NewSalaryCalculator()

	tests := []struct {
		name     string
		salaries []models.Salary
		expected []float64
	}{
		{
			name: "December bonus adjustment",
			salaries: []models.Salary{
				{Month: 12, Salary: 2000},
			},
			expected: []float64{2200}, // 2000 * 1.10
		},
		{
			name: "Summer months adjustment - June",
			salaries: []models.Salary{
				{Month: 6, Salary: 2000},
			},
			expected: []float64{1900}, // 2000 * 0.95
		},
		{
			name: "Summer months adjustment - July",
			salaries: []models.Salary{
				{Month: 7, Salary: 2000},
			},
			expected: []float64{1900}, // 2000 * 0.95
		},
		{
			name: "Summer months adjustment - August",
			salaries: []models.Salary{
				{Month: 8, Salary: 2000},
			},
			expected: []float64{1900}, // 2000 * 0.95
		},
		{
			name: "Regular months - no adjustment",
			salaries: []models.Salary{
				{Month: 1, Salary: 2000},
				{Month: 3, Salary: 1500},
				{Month: 9, Salary: 1800},
			},
			expected: []float64{2000, 1500, 1800},
		},
		{
			name: "Mixed seasonal adjustments",
			salaries: []models.Salary{
				{Month: 12, Salary: 1800}, // +10%
				{Month: 6, Salary: 1900},  // -5%
				{Month: 1, Salary: 2000},  // no adjustment
			},
			expected: []float64{1980, 1805, 2000}, // 1800*1.10, 1900*0.95, 2000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.applySeasonalAdjustments(tt.salaries)
			
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if abs(result[i]-expected) > 0.01 { // Allow small floating point differences
					t.Errorf("Expected salary[%d] = %.2f, got %.2f", i, expected, result[i])
				}
			}
		})
	}
}

func TestSalaryCalculator_ApplyTaxDeduction(t *testing.T) {
	calculator := NewSalaryCalculator()

	tests := []struct {
		name     string
		salaries []float64
		total    float64
		expected []float64
	}{
		{
			name:     "Below tax threshold - no deduction",
			salaries: []float64{2000, 2500, 3000},
			total:    7500,
			expected: []float64{2000, 2500, 3000},
		},
		{
			name:     "Above tax threshold - apply 7% deduction",
			salaries: []float64{3000, 3100, 3200, 3300},
			total:    12600,
			expected: []float64{2790, 2883, 2976, 3069}, // each * 0.93
		},
		{
			name:     "Exactly at threshold - no deduction",
			salaries: []float64{5000, 5000},
			total:    10000,
			expected: []float64{5000, 5000},
		},
		{
			name:     "Just above threshold - apply deduction",
			salaries: []float64{5000, 5001},
			total:    10001,
			expected: []float64{4650, 4650.93}, // each * 0.93
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.applyTaxDeduction(tt.salaries, tt.total)
			
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if abs(result[i]-expected) > 0.01 { // Allow small floating point differences
					t.Errorf("Expected salary[%d] = %.2f, got %.2f", i, expected, result[i])
				}
			}
		})
	}
}

func TestSalaryCalculator_DetermineStatus(t *testing.T) {
	calculator := NewSalaryCalculator()

	tests := []struct {
		name           string
		averageSalary  float64
		expectedStatus string
	}{
		{
			name:           "Average above 2000 - GREEN",
			averageSalary:  2500.75,
			expectedStatus: "GREEN",
		},
		{
			name:           "Average exactly 2000 - ORANGE",
			averageSalary:  2000.0,
			expectedStatus: "ORANGE",
		},
		{
			name:           "Average below 2000 - RED",
			averageSalary:  1999.99,
			expectedStatus: "RED",
		},
		{
			name:           "Very low average - RED",
			averageSalary:  500.0,
			expectedStatus: "RED",
		},
		{
			name:           "Very high average - GREEN",
			averageSalary:  10000.0,
			expectedStatus: "GREEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.determineStatus(tt.averageSalary)
			if result != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s for average %.2f", tt.expectedStatus, result, tt.averageSalary)
			}
		})
	}
}

func TestSalaryCalculator_CalculateEmployeeStatus_CompleteFlow(t *testing.T) {
	calculator := NewSalaryCalculator()

	tests := []struct {
		name            string
		salaries        []models.Salary
		expectedStatus  string
		expectedAverage float64
		expectedHighest float64
		expectedSum     float64
	}{
		{
			name: "NAT1001 scenario - RED status with summer adjustment",
			salaries: []models.Salary{
				{Month: 1, Salary: 1200},
				{Month: 2, Salary: 1300},
				{Month: 3, Salary: 1400},
				{Month: 5, Salary: 1500},
				{Month: 6, Salary: 1600}, // Summer: 1600 * 0.95 = 1520
			},
			expectedStatus:  "RED",
			expectedAverage: 1384,     // (1200+1300+1400+1500+1520)/5
			expectedHighest: 1520,     // Adjusted June salary
			expectedSum:     6920,     // Sum of adjusted salaries
		},
		{
			name: "NAT1002 scenario - ORANGE status",
			salaries: []models.Salary{
				{Month: 1, Salary: 1900},
				{Month: 2, Salary: 2000},
				{Month: 3, Salary: 2100},
			},
			expectedStatus:  "ORANGE",
			expectedAverage: 2000,
			expectedHighest: 2100,
			expectedSum:     6000,
		},
		{
			name: "NAT1004 scenario - GREEN status with tax deduction",
			salaries: []models.Salary{
				{Month: 1, Salary: 2500},
				{Month: 2, Salary: 2600},
				{Month: 3, Salary: 2700},
				{Month: 4, Salary: 2800},
			},
			// Total before tax: 10600 > 10000, so 7% tax applies
			// After tax: 2325, 2418, 2511, 2604
			expectedStatus:  "GREEN",
			expectedAverage: 2464.5, // (2325+2418+2511+2604)/4
			expectedHighest: 2604,   // 2800 * 0.93
			expectedSum:     9858,   // Sum after tax
		},
		{
			name: "NAT1005 scenario - HIGH salary with tax deduction",
			salaries: []models.Salary{
				{Month: 1, Salary: 3000},
				{Month: 2, Salary: 3100},
				{Month: 3, Salary: 3200},
				{Month: 4, Salary: 3300},
			},
			expectedStatus:  "GREEN",
			expectedAverage: 2929.5, // (2790+2883+2976+3069)/4 after 7% tax
			expectedHighest: 3069,   // 3300 * 0.93
			expectedSum:     11718,  // Sum after tax deduction
		},
		{
			name: "December bonus scenario",
			salaries: []models.Salary{
				{Month: 12, Salary: 2000}, // 2000 * 1.10 = 2200
				{Month: 1, Salary: 1800},
				{Month: 2, Salary: 1900},
			},
			expectedStatus:  "RED",
			expectedAverage: 1966.67, // (2200+1800+1900)/3
			expectedHighest: 2200,    // December bonus adjusted
			expectedSum:     5900,
		},
		{
			name: "Summer months scenario",
			salaries: []models.Salary{
				{Month: 6, Salary: 2000}, // 2000 * 0.95 = 1900
				{Month: 7, Salary: 2100}, // 2100 * 0.95 = 1995
				{Month: 8, Salary: 2200}, // 2200 * 0.95 = 2090
			},
			expectedStatus:  "RED",
			expectedAverage: 1995,    // (1900+1995+2090)/3
			expectedHighest: 2090,    // August adjusted
			expectedSum:     5985,
		},
		{
			name: "Border case - exactly 2000 average",
			salaries: []models.Salary{
				{Month: 1, Salary: 1950},
				{Month: 2, Salary: 2000},
				{Month: 3, Salary: 2050},
			},
			expectedStatus:  "ORANGE",
			expectedAverage: 2000,
			expectedHighest: 2050,
			expectedSum:     6000,
		},
		{
			name: "Empty salaries",
			salaries:        []models.Salary{},
			expectedStatus:  "RED",
			expectedAverage: 0,
			expectedHighest: 0,
			expectedSum:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateEmployeeStatus(tt.salaries)

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}

			if abs(result.AverageSalary-tt.expectedAverage) > 0.01 {
				t.Errorf("Expected average %.2f, got %.2f", tt.expectedAverage, result.AverageSalary)
			}

			if abs(result.HighestSalary-tt.expectedHighest) > 0.01 {
				t.Errorf("Expected highest %.2f, got %.2f", tt.expectedHighest, result.HighestSalary)
			}

			if abs(result.SumOfSalaries-tt.expectedSum) > 0.01 {
				t.Errorf("Expected sum %.2f, got %.2f", tt.expectedSum, result.SumOfSalaries)
			}
		})
	}
}

func TestSalaryCalculator_EdgeCases(t *testing.T) {
	calculator := NewSalaryCalculator()

	t.Run("Single salary record", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 2500},
		}
		result := calculator.CalculateEmployeeStatus(salaries)
		
		if result.Status != "GREEN" {
			t.Errorf("Expected GREEN status for single high salary, got %s", result.Status)
		}
		if result.AverageSalary != 2500 {
			t.Errorf("Expected average 2500, got %.2f", result.AverageSalary)
		}
	})

	t.Run("All summer months", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 6, Salary: 2500}, // 2375
			{Month: 7, Salary: 2500}, // 2375
			{Month: 8, Salary: 2500}, // 2375
		}
		result := calculator.CalculateEmployeeStatus(salaries)
		
		if result.Status != "GREEN" {
			t.Errorf("Expected GREEN status after summer adjustments, got %s", result.Status)
		}
		expectedAverage := 2375.0 // 2500 * 0.95
		if abs(result.AverageSalary-expectedAverage) > 0.01 {
			t.Errorf("Expected average %.2f, got %.2f", expectedAverage, result.AverageSalary)
		}
	})

	t.Run("Mixed adjustments with tax deduction", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 12, Salary: 3000}, // 3300 before tax, 3069 after tax
			{Month: 6, Salary: 3000},  // 2850 before tax, 2650.5 after tax
			{Month: 1, Salary: 3000},  // 3000 before tax, 2790 after tax
			{Month: 2, Salary: 3000},  // 3000 before tax, 2790 after tax
		}
		// Total before tax: 3300 + 2850 + 3000 + 3000 = 12150 > 10000
		// After tax: each * 0.93
		result := calculator.CalculateEmployeeStatus(salaries)
		
		if result.Status != "GREEN" {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}
		
		// Expected: (3069 + 2650.5 + 2790 + 2790) / 4 = 2824.875
		expectedAverage := 2824.875
		if abs(result.AverageSalary-expectedAverage) > 0.01 {
			t.Errorf("Expected average %.3f, got %.3f", expectedAverage, result.AverageSalary)
		}
	})
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}