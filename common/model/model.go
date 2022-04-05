package model

import (
	"time"
)

type Task struct {
	TaskId            string `gorm:"primaryKey"`
	UserId            string
	Status            int
	StatusDescription string
	Title             string
	Content           string `gorm:"type:longtext"`
	StartTime         time.Time
}

type User struct {
	UserId   string `gorm:"primaryKey"`
	Username string
	Password string
	Phone    string
}
