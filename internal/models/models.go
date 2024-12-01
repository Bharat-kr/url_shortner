package models

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	id           uint   `gorm:"primaryKey;default:auto_random()"`
	original_url string `gorm:"not null"`
	short_url    string `gorm:"not null"`
}
