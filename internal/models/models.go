package models

import "time"

type NewURL struct {
	FullURL string `json:"full_url" binding:"required,url"`
}

type URL struct {
	Code       string    `json:"code"`
	FullURL    string    `json:"full_url"`
	IsActive   int       `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int       `json:"click_count"`
}

type UpdateURL struct {
	FullURL  string `json:"full_url"`
	IsActive int    `json:"is_active"`
}
