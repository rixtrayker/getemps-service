package models

import "time"

type EmployeeInfo struct {
	ID             int64         `json:"id"`
	Username       string        `json:"username"`
	NationalNumber string        `json:"nationalNumber"`
	Email          string        `json:"email"`
	Phone          string        `json:"phone"`
	IsActive       bool          `json:"isActive"`
	SalaryDetails  SalaryDetails `json:"salaryDetails"`
	Status         string        `json:"status"`
	LastUpdated    time.Time     `json:"lastUpdated"`
}

type SalaryDetails struct {
	AverageSalary float64 `json:"averageSalary"`
	HighestSalary float64 `json:"highestSalary"`
	SumOfSalaries float64 `json:"sumOfSalaries"`
}