package validator

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yourusername/getemps-service/internal/models"
)

func ValidateEmployeeRequest(req models.EmployeeRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.NationalNumber,
			validation.Required.Error("National number is required"),
			validation.Length(3, 50).Error("Invalid national number format"),
		),
	)
}

func IsValidNationalNumber(nationalNumber string) bool {
	if len(nationalNumber) < 3 || len(nationalNumber) > 50 {
		return false
	}
	
	// Additional validation logic can be added here
	// For now, we just check length and non-empty
	return nationalNumber != ""
}