package service

import (
	"testing"

	"github.com/rixtrayker/getemps-service/internal/models"
)

// ============================================
// Basic Seasonal Adjustments Tests
// ============================================

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
		{
			name: "Multiple December bonuses",
			salaries: []models.Salary{
				{Month: 12, Year: 2023, Salary: 2000},
				{Month: 12, Year: 2024, Salary: 2000},
				{Month: 12, Year: 2025, Salary: 2000},
			},
			expected: []float64{2200, 2200, 2200}, // All get +10%
		},
		{
			name: "All summer months",
			salaries: []models.Salary{
				{Month: 6, Salary: 2500},
				{Month: 7, Salary: 2500},
				{Month: 8, Salary: 2500},
			},
			expected: []float64{2375, 2375, 2375}, // All get -5%
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

// ============================================
// Tax Deduction Tests
// ============================================

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
		{
			name:     "Exactly 10000.01 - apply deduction",
			salaries: []float64{2500, 2500, 2500, 2500.01},
			total:    10000.01,
			expected: []float64{2325, 2325, 2325, 2325.0093}, // each * 0.93
		},
		{
			name:     "Very high salaries",
			salaries: []float64{20000, 21000, 22000},
			total:    63000,
			expected: []float64{18600, 19530, 20460}, // each * 0.93
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

// ============================================
// Status Determination Tests
// ============================================

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
			expectedStatus: StatusGreen,
		},
		{
			name:           "Average exactly 2000 - ORANGE",
			averageSalary:  2000.0,
			expectedStatus: StatusOrange,
		},
		{
			name:           "Average below 2000 - RED",
			averageSalary:  1999.99,
			expectedStatus: StatusRed,
		},
		{
			name:           "Very low average - RED",
			averageSalary:  500.0,
			expectedStatus: StatusRed,
		},
		{
			name:           "Very high average - GREEN",
			averageSalary:  10000.0,
			expectedStatus: StatusGreen,
		},
		{
			name:           "Just above 2000 - GREEN",
			averageSalary:  2000.01,
			expectedStatus: StatusGreen,
		},
		{
			name:           "Just below 2000 - RED",
			averageSalary:  1999.999,
			expectedStatus: StatusRed,
		},
		{
			name:           "Zero average - RED",
			averageSalary:  0,
			expectedStatus: StatusRed,
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

// ============================================
// Complete Flow Tests - Real Scenarios
// ============================================

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
			expectedStatus:  StatusRed,
			expectedAverage: 1384, // (1200+1300+1400+1500+1520)/5
			expectedHighest: 1520, // Adjusted June salary
			expectedSum:     6920, // Sum of adjusted salaries
		},
		{
			name: "Real NAT1002 scenario - RED status (low salaries)",
			salaries: []models.Salary{
				{Month: 1, Salary: 900},
				{Month: 2, Salary: 950},
				{Month: 3, Salary: 980},
				{Month: 4, Salary: 1100},
				{Month: 5, Salary: 1150},
			},
			expectedStatus:  StatusRed,
			expectedAverage: 1016,
			expectedHighest: 1150,
			expectedSum:     5080,
		},
		{
			name: "ORANGE status scenario - exactly 2000 average",
			salaries: []models.Salary{
				{Month: 1, Salary: 1900},
				{Month: 2, Salary: 2000},
				{Month: 3, Salary: 2100},
			},
			expectedStatus:  StatusOrange,
			expectedAverage: 2000,
			expectedHighest: 2100,
			expectedSum:     6000,
		},
		{
			name: "NAT1004 scenario - RED status (corrected expectation)",
			salaries: []models.Salary{
				{Month: 1, Salary: 2000},
				{Month: 2, Salary: 2050},
				{Month: 3, Salary: 2100},
				{Month: 4, Salary: 2200},
				{Month: 5, Salary: 2300},
			},
			// Total: 10650 > 10000, tax applies
			// After 7% tax: 1860, 1906.5, 1953, 2046, 2139
			expectedStatus:  StatusRed, // Average after tax is 1980.9 < 2000
			expectedAverage: 1980.9,    // Average after tax
			expectedHighest: 2139,      // 2300 * 0.93
			expectedSum:     9904.5,    // Sum after tax
		},
		{
			name: "HIGH salary with tax deduction",
			salaries: []models.Salary{
				{Month: 1, Salary: 3000},
				{Month: 2, Salary: 3100},
				{Month: 3, Salary: 3200},
				{Month: 4, Salary: 3300},
			},
			// Total: 12600 > 10000, tax applies
			expectedStatus:  StatusGreen,
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
			expectedStatus:  StatusRed,
			expectedAverage: 1966.67, // (2200+1800+1900)/3
			expectedHighest: 2200,    // December bonus adjusted
			expectedSum:     5900,
		},
		{
			name: "All summer months scenario",
			salaries: []models.Salary{
				{Month: 6, Salary: 2000}, // 2000 * 0.95 = 1900
				{Month: 7, Salary: 2100}, // 2100 * 0.95 = 1995
				{Month: 8, Salary: 2200}, // 2200 * 0.95 = 2090
			},
			expectedStatus:  StatusRed,
			expectedAverage: 1995, // (1900+1995+2090)/3
			expectedHighest: 2090, // August adjusted
			expectedSum:     5985,
		},
		{
			name: "Border case - exactly 2000 average",
			salaries: []models.Salary{
				{Month: 1, Salary: 1950},
				{Month: 2, Salary: 2000},
				{Month: 3, Salary: 2050},
			},
			expectedStatus:  StatusOrange,
			expectedAverage: 2000,
			expectedHighest: 2050,
			expectedSum:     6000,
		},
		{
			name: "December and summer months together",
			salaries: []models.Salary{
				{Month: 6, Salary: 2000},  // -5% = 1900
				{Month: 12, Salary: 2000}, // +10% = 2200
				{Month: 1, Salary: 2000},  // No change
			},
			expectedStatus:  StatusGreen, // Average 2033.33 > 2000
			expectedAverage: 2033.33,     // (1900+2200+2000)/3
			expectedHighest: 2200,
			expectedSum:     6100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateEmployeeStatus(tt.salaries)

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}

			if abs(result.AverageSalary-tt.expectedAverage) > 0.5 {
				t.Errorf("Expected average %.2f, got %.2f", tt.expectedAverage, result.AverageSalary)
			}

			if abs(result.HighestSalary-tt.expectedHighest) > 0.5 {
				t.Errorf("Expected highest %.2f, got %.2f", tt.expectedHighest, result.HighestSalary)
			}

			if abs(result.SumOfSalaries-tt.expectedSum) > 0.5 {
				t.Errorf("Expected sum %.2f, got %.2f", tt.expectedSum, result.SumOfSalaries)
			}
		})
	}
}

// ============================================
// Critical Edge Cases Tests
// ============================================

func TestSalaryCalculator_CriticalEdgeCases(t *testing.T) {
	calculator := NewSalaryCalculator()

	t.Run("Tax threshold exactly 10000 - no tax applied", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 2500},
			{Month: 2, Salary: 2500},
			{Month: 3, Salary: 2500},
			{Month: 4, Salary: 2500},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		// Total is exactly 10000, should NOT apply tax (requirement: > 10000)
		if abs(result.SumOfSalaries-10000) > 0.01 {
			t.Errorf("Expected sum 10000 (no tax), got %.2f", result.SumOfSalaries)
		}
		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}
	})

	t.Run("Tax threshold 10000.01 - tax applied", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 2500},
			{Month: 2, Salary: 2500},
			{Month: 3, Salary: 2500},
			{Month: 4, Salary: 2500.01},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		// Total is 10000.01 > 10000, should apply 7% tax
		expectedSum := 10000.01 * 0.93 // 9300.0093
		if abs(result.SumOfSalaries-expectedSum) > 0.01 {
			t.Errorf("Expected sum %.2f (with tax), got %.2f", expectedSum, result.SumOfSalaries)
		}
	})

	t.Run("Status boundary - exactly 2000 after all adjustments", func(t *testing.T) {
		// This tests floating point precision at the ORANGE boundary
		salaries := []models.Salary{
			{Month: 1, Salary: 2000.00},
			{Month: 2, Salary: 2000.00},
			{Month: 3, Salary: 2000.00},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusOrange {
			t.Errorf("Expected ORANGE for exactly 2000 average, got %s", result.Status)
		}
		if abs(result.AverageSalary-2000) > 0.01 {
			t.Errorf("Expected average 2000.00, got %.2f", result.AverageSalary)
		}
	})

	t.Run("Just below 2000 average - RED", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 1999.99},
			{Month: 2, Salary: 2000.00},
			{Month: 3, Salary: 2000.00},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusRed {
			t.Errorf("Expected RED for average < 2000, got %s", result.Status)
		}
	})

	t.Run("Just above 2000 average - GREEN", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 2000.01},
			{Month: 2, Salary: 2000.00},
			{Month: 3, Salary: 2000.00},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN for average > 2000, got %s", result.Status)
		}
	})

	t.Run("Summer adjustment keeps total below 10000 - no tax", func(t *testing.T) {
		// Summer deduction keeps total below 10000
		salaries := []models.Salary{
			{Month: 6, Salary: 3400}, // Summer: 3400 * 0.95 = 3230
			{Month: 7, Salary: 3400}, // Summer: 3400 * 0.95 = 3230
			{Month: 8, Salary: 3400}, // Summer: 3400 * 0.95 = 3230
		}
		// Total after seasonal: 9690 < 10000, no tax
		result := calculator.CalculateEmployeeStatus(salaries)

		expectedSum := 9690.0
		if abs(result.SumOfSalaries-expectedSum) > 0.01 {
			t.Errorf("Expected sum %.2f (no tax due to seasonal), got %.2f", expectedSum, result.SumOfSalaries)
		}
		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}
	})

	t.Run("December bonus pushes total over 10000 - tax applied", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 12, Salary: 3400}, // +10% = 3740
			{Month: 1, Salary: 3200},
			{Month: 2, Salary: 3200},
		}
		// Total after seasonal: 10140 > 10000, tax applies
		// After tax: 3478.2, 2976, 2976 = 9430.2
		result := calculator.CalculateEmployeeStatus(salaries)

		expectedSum := 9430.2
		if abs(result.SumOfSalaries-expectedSum) > 0.5 {
			t.Errorf("Expected sum %.2f (with tax after seasonal), got %.2f", expectedSum, result.SumOfSalaries)
		}
	})

	t.Run("Mixed adjustments with high tax impact", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 12, Salary: 3000}, // +10% = 3300
			{Month: 6, Salary: 3000},  // -5% = 2850
			{Month: 1, Salary: 3000},  // No change = 3000
			{Month: 2, Salary: 3000},  // No change = 3000
		}
		// Total before tax: 3300 + 2850 + 3000 + 3000 = 12150 > 10000
		// After tax: 3069, 2650.5, 2790, 2790 = 11299.5
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}

		// Expected average: (3069 + 2650.5 + 2790 + 2790) / 4 = 2824.875
		expectedAverage := 2824.875
		if abs(result.AverageSalary-expectedAverage) > 1.0 {
			t.Errorf("Expected average %.3f, got %.3f", expectedAverage, result.AverageSalary)
		}
	})
}

// ============================================
// Additional Edge Cases Tests
// ============================================

func TestSalaryCalculator_AdditionalEdgeCases(t *testing.T) {
	calculator := NewSalaryCalculator()

	t.Run("Single salary record", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 2500},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status for single high salary, got %s", result.Status)
		}
		if abs(result.AverageSalary-2500) > 0.01 {
			t.Errorf("Expected average 2500, got %.2f", result.AverageSalary)
		}
	})

	t.Run("All summer months with tax", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 6, Salary: 3600}, // 3420
			{Month: 7, Salary: 3600}, // 3420
			{Month: 8, Salary: 3600}, // 3420
		}
		// Total: 10260 > 10000, tax applies
		// After tax: 3180.6, 3180.6, 3180.6
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status after summer adjustments and tax, got %s", result.Status)
		}
	})

	t.Run("Empty salaries slice", func(t *testing.T) {
		salaries := []models.Salary{}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusRed {
			t.Errorf("Expected RED status for empty salaries, got %s", result.Status)
		}
		if result.AverageSalary != 0 {
			t.Errorf("Expected average 0, got %.2f", result.AverageSalary)
		}
	})

	t.Run("Zero salary values", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 0},
			{Month: 2, Salary: 0},
			{Month: 3, Salary: 0},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusRed {
			t.Errorf("Expected RED status for zero salaries, got %s", result.Status)
		}
		if result.AverageSalary != 0 {
			t.Errorf("Expected average 0, got %.2f", result.AverageSalary)
		}
	})

	t.Run("Very large salary values", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 50000},
			{Month: 2, Salary: 51000},
			{Month: 3, Salary: 52000},
		}
		// Total: 153000 > 10000, tax applies
		// After tax: 46500, 47430, 48360 = 142290
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status for very high salaries, got %s", result.Status)
		}
	})

	t.Run("All December months (multiple years)", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 12, Year: 2023, Salary: 2000}, // +10%
			{Month: 12, Year: 2024, Salary: 2000}, // +10%
			{Month: 12, Year: 2025, Salary: 2000}, // +10%
		}
		// All get +10%: 2200, 2200, 2200 = 6600
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}
		expectedAverage := 2200.0
		if abs(result.AverageSalary-expectedAverage) > 0.01 {
			t.Errorf("Expected average %.2f, got %.2f", expectedAverage, result.AverageSalary)
		}
	})

	t.Run("December and all summer months", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 6, Salary: 2000},  // -5% = 1900
			{Month: 7, Salary: 2000},  // -5% = 1900
			{Month: 8, Salary: 2000},  // -5% = 1900
			{Month: 12, Salary: 2000}, // +10% = 2200
		}
		// Total: 7900 < 10000, no tax
		// Average: 1975
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusRed {
			t.Errorf("Expected RED status, got %s", result.Status)
		}
		expectedAverage := 1975.0
		if abs(result.AverageSalary-expectedAverage) > 0.01 {
			t.Errorf("Expected average %.2f, got %.2f", expectedAverage, result.AverageSalary)
		}
	})

	t.Run("Salary with decimal precision", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 1, Salary: 1999.99},
			{Month: 2, Salary: 2000.00},
			{Month: 3, Salary: 2000.01},
		}
		result := calculator.CalculateEmployeeStatus(salaries)

		// Average should be exactly 2000.00
		if result.Status != StatusOrange {
			t.Errorf("Expected ORANGE status, got %s", result.Status)
		}
	})

	t.Run("Complex scenario - NAT1008 high earner", func(t *testing.T) {
		salaries := []models.Salary{
			{Month: 10, Salary: 2200},
			{Month: 11, Salary: 2300},
			{Month: 12, Salary: 2400}, // +10% = 2640
			{Month: 1, Salary: 2500},
			{Month: 2, Salary: 2600},
			{Month: 3, Salary: 2800},
		}
		// Total before tax: 2200 + 2300 + 2640 + 2500 + 2600 + 2800 = 15040 > 10000
		// After 7% tax all values * 0.93
		// After tax: 2046, 2139, 2455.2, 2325, 2418, 2604 = 13987.2
		result := calculator.CalculateEmployeeStatus(salaries)

		if result.Status != StatusGreen {
			t.Errorf("Expected GREEN status, got %s", result.Status)
		}

		// Average: 13987.2 / 6 = 2331.2
		expectedAverage := 2331.2
		if abs(result.AverageSalary-expectedAverage) > 1.0 {
			t.Errorf("Expected average %.2f, got %.2f", expectedAverage, result.AverageSalary)
		}
	})
}

// ============================================
// Helper Functions
// ============================================

// Helper function for floating point comparison.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
