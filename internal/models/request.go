package models

type EmployeeRequest struct {
	NationalNumber string `json:"NationalNumber"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}