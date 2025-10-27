package service

import (
	"github.com/rixtrayker/getemps-service/internal/models"
)

type SalaryCalculator struct{}

func NewSalaryCalculator() *SalaryCalculator {
	return &SalaryCalculator{}
}

type SalaryCalculationResult struct {
	AverageSalary float64
	HighestSalary float64
	SumOfSalaries float64
	Status        string
}

func (sc *SalaryCalculator) CalculateEmployeeStatus(salaries []models.Salary) *SalaryCalculationResult {
	if len(salaries) == 0 {
		return &SalaryCalculationResult{
			AverageSalary: 0,
			HighestSalary: 0,
			SumOfSalaries: 0,
			Status:        "RED",
		}
	}

	// Step 1: Apply seasonal adjustments
	adjustedSalaries := sc.applySeasonalAdjustments(salaries)

	// Step 2: Calculate total before tax
	totalBeforeTax := sc.sumSalaries(adjustedSalaries)

	// Step 3: Apply tax deduction if needed
	finalSalaries := sc.applyTaxDeduction(adjustedSalaries, totalBeforeTax)

	// Step 4: Calculate final statistics
	return sc.calculateFinalStats(finalSalaries)
}

func (sc *SalaryCalculator) applySeasonalAdjustments(salaries []models.Salary) []float64 {
	adjusted := make([]float64, len(salaries))

	for i, salary := range salaries {
		adjustedSalary := salary.Salary

		// December: +10% holiday bonus
		if salary.Month == 12 {
			adjustedSalary = adjustedSalary * 1.10
		}

		// Summer months (June, July, August): -5% deduction
		if salary.Month == 6 || salary.Month == 7 || salary.Month == 8 {
			adjustedSalary = adjustedSalary * 0.95
		}

		adjusted[i] = adjustedSalary
	}

	return adjusted
}

func (sc *SalaryCalculator) sumSalaries(salaries []float64) float64 {
	total := 0.0
	for _, salary := range salaries {
		total += salary
	}
	return total
}

func (sc *SalaryCalculator) applyTaxDeduction(salaries []float64, total float64) []float64 {
	// If total salary > 10,000: apply 7% tax deduction to each salary
	if total > 10000 {
		taxedSalaries := make([]float64, len(salaries))
		for i, salary := range salaries {
			taxedSalaries[i] = salary * 0.93 // 7% tax deduction
		}
		return taxedSalaries
	}

	return salaries
}

func (sc *SalaryCalculator) calculateFinalStats(salaries []float64) *SalaryCalculationResult {
	if len(salaries) == 0 {
		return &SalaryCalculationResult{
			AverageSalary: 0,
			HighestSalary: 0,
			SumOfSalaries: 0,
			Status:        "RED",
		}
	}

	// Calculate sum
	sum := sc.sumSalaries(salaries)

	// Calculate average
	average := sum / float64(len(salaries))

	// Find highest
	highest := salaries[0]
	for _, salary := range salaries {
		if salary > highest {
			highest = salary
		}
	}

	// Determine status based on average
	status := sc.determineStatus(average)

	return &SalaryCalculationResult{
		AverageSalary: average,
		HighestSalary: highest,
		SumOfSalaries: sum,
		Status:        status,
	}
}

func (sc *SalaryCalculator) determineStatus(averageSalary float64) string {
	if averageSalary > 2000 {
		return "GREEN"
	} else if averageSalary == 2000 {
		return "ORANGE"
	} else {
		return "RED"
	}
}