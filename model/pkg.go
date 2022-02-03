package model

import (
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/datatypes"
)

var (
	Settings = &settings.Settings
)

//
// Field (data) types.
type JSON = datatypes.JSON

//
// All builds all models.
// Models are enumerated such that each are listed after
// all the other models on which they may depend.
func All() []interface{} {
	return []interface{}{
		Setting{},
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
		Bucket{},
		Dependency{},
		Review{},
		Identity{},
		Task{},
		TaskReport{},
		Proxy{},
	}
}
