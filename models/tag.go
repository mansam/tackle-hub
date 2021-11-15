package models

type Tag struct {
	Resource
	Name      string `json:"name" gorm:"notnull" binding:"required"`
	TagTypeID string `json:"tag_type_id" gorm:"notnull" binding:"required"`
	TagType   TagType
}
