package models

type Application struct {
	Resource
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Comments          string        `json:"comments"`
	BusinessServiceID string        `json:"business_service_id" gorm:"notnull" binding:"required"`
	DependsOn         []Application `json:"depends_on" gorm:"many2many:application_dependencies"`
	Tags              []Tag         `json:"tags" gorm:"many2many:application_tags"`
	BusinessService   BusinessService
}
