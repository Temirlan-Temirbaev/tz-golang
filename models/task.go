package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null"`
	TwitterID string `json:"twitter_id" gorm:"unique;not null"`
	Users     []User `gorm:"many2many:user_tasks;"`
}
