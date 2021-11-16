package models

type BusinessService struct {
	Resource
	Name        string        `json:"name" gorm:"notnull,unique" validate:"required"`
	Description string        `json:"description"`
	Users       []Stakeholder `json:"users" gorm:"many2many:user_business_services"`
}
