package models

type Stakeholder struct {
	Resource
	Email            string `json:"email" gorm:"notnull" binding:"required,email"`
	DisplayName      string `json:"display_name" gorm:"notnull" binding:"required"`
	JobFunctionID    string `json:"job_function_id" gorm:"notnull" binding:"required"`
	JobFunction      JobFunction
	BusinessServices []BusinessService `json:"business_services" gorm:"many2many:user_business_services"`
	Groups           []Group           `json:"groups" gorm:"many2many:user_groups"`
}
