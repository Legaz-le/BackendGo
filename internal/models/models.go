package models

import "time"

type Bookmark struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Job struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	SalaryRange string    `json:"salary_range"`
	CreatedAt   time.Time `json:"created_at"`
}
