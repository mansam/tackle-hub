package models

type SourceRepo struct {
	Resource
	Type          string `json:"name" gorm:"notnull" binding:"required" validate:"oneof= git svn"`
	URL           string `json:"url" gorm:"notnull" binding:"required"`
	Branch        string `json:"branch" gorm:"notnull" binding:"required"`
	ApplicationID string `json:"application_id"`
	Application   Application
}
