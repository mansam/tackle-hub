package model

import "time"

//
// Model Base model.
type Model struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreateUser string    `json:"createUser"`
	UpdateUser string    `json:"updateUser"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
}
