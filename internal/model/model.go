package model

import "time"

type Item struct {
	ID         int        `json:"id"`
	Type       string     `json:"type"`
	Amount     int        `json:"amount"`
	Date       string     `json:"date"`
	Category   string     `json:"category"`
	CreatedAt  time.Time  `json:"created_at"`
	Aggregated Aggregated `json:"aggregated_data,omitempty"`
}

type Aggregated struct {
	Sum           int     `json:"sum"`
	Average       float64 `json:"average"`
	Count         int     `json:"count"`
	Median        float64 `json:"median"`
	Percentile_90 float64 `json:"percentile_90"`
}
