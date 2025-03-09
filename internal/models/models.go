package models

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	OriginalUrl string `gorm:"not null"`
	ShortUrl    string `gorm:"not null"`
}
