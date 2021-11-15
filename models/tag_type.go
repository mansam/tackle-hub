package models

type TagType struct {
	Resource
	Name   string `json:"name" gorm:"notnull" binding:"required"`
	Rank   uint   `json:"rank"`
	Colour string `json:"colour"`
	Tags   []Tag  `json:"tags"`
}
