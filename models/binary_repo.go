package models

type BinaryRepo struct {
	Resource
	Type          string `json:"name" gorm:"notnull" binding:"required" validate:"oneof=mvn"`
	URL           string `json:"url" gorm:"notnull" binding:"required"`
	Group         string `json:"group" gorm:"notnull" binding:"required"`
	Artifact      string `json:"artifact" gorm:"notnull" binding:"required"`
	Version       string `json:"version" gorm:"notnull" binding:"required"`
	ApplicationID string `json:"application_id"`
	Application   Application
}
