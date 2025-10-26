package models

import "time"

type Salary struct {
	ID        int64     `json:"id" db:"id"`
	Year      int       `json:"year" db:"year"`
	Month     int       `json:"month" db:"month"`
	Salary    float64   `json:"salary" db:"salary"`
	UserID    int64     `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}