package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID   string `json:"telegram_id" gorm:"primaryKey;unique;not null"`
	AccessToken  string `json:"-" gorm:"not null"`
	AccessSecret string `json:"-" gorm:"not null"`
	Tasks        []Task `gorm:"many2many:user_tasks;"`
}
