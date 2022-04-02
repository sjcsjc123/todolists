package model

import (
	"time"
)

type Task struct {
	TaskId            int `gorm:"primaryKey"`
	UserId            int
	Status            int
	StatusDescription string
	Title             string
	Content           string `gorm:"type:longtext"`
	StartTime         time.Time
}

type User struct {
	UserId   int `gorm:"primaryKey"`
	Username string
	Password string
	Phone    string
}
