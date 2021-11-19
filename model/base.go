package model

//
// Model Base model.
type Model struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	CreateUser string `json:"createUser"`
	UpdateUser string `json:"updateUser"`
}
