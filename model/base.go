package model

import (
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/datatypes"
	"time"
)

var (
	Settings = &settings.Settings
)

//
// Field (data) types.
type JSON = datatypes.JSON

//
// Model Base model.
type Model struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreateUser string    `json:"createUser"`
	UpdateUser string    `json:"updateUser"`
	CreateTime time.Time `json:"createTime" gorm:"column:createTime;autoCreateTime"`
}

//
// Seeded model.
type Seeded struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

//
// All builds all models.
// Models are enumerated such that each are listed after
// all the other models on which they may depend.
func All() []interface{} {
	return []interface{}{
		ImportSummary{},
		Import{},
		ImportTag{},
		JobFunction{},
		TagType{},
		Tag{},
		StakeholderGroup{},
		Stakeholder{},
		BusinessService{},
		Application{},
		Dependency{},
		Review{},
		Repository{},
		Identity{},
		Task{},
		TaskReport{},
	}
}
