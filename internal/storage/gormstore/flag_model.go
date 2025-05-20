package gormstore

import "time"

type Flag struct {
	FlagID    string `gorm:"primaryKey"`
	PageID    string
	Status    string
	Body      string
	UserID    string
	UpdatedAt time.Time
}
