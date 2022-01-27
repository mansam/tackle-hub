package model

import (
	"time"
)

//
// Model Base model.
type Model struct {
	ID         uint `gorm:"primaryKey"`
	CreateUser string
	UpdateUser string
	CreateTime time.Time `gorm:"column:createTime;autoCreateTime"`
}

//
// Seeded model.
type Seeded struct {
	ID uint `gorm:"primaryKey"`
}
