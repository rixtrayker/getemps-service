package models

import "time"

type User struct {
	ID             int64     `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	NationalNumber string    `json:"nationalNumber" db:"national_number"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	IsActive       bool      `json:"isActive" db:"is_active"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
