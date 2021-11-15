package models

type JobFunction struct {
	Resource
	Name string `json:"name" gorm:"notnull,unique"`
}
