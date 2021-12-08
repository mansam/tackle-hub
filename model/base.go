package model

import "time"

//
// Model Base model.
type Model struct {
	CreateTime time.Time `json:"createTime" gorm:"column:createTime; autoCreateTime"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreateUser string    `json:"createUser"`
	UpdateUser string    `json:"updateUser"`
}
