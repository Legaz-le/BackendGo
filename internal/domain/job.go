package domain

import (
	"time"
)


type Job struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	SalaryMin   int       `json:"min_salary"`
	SalaryMax   int       `json:"max_salary"`
	CreatedAt   time.Time `json:"created_at"`
}