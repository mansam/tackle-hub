package models

type Group struct {
	Resource
	Name        string `json:"name" gorm:"unique,index" validate:"required,alphanum,min=6,max=32"`
	Description string `json:"description"`
	Members     []User `json:"members" gorm:"many2many:user_groups"`
}
