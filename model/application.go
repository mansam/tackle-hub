package model

type Application struct {
	Model
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Review            *Review          `json:"review"`
	Comments          string           `json:"comments"`
	Tags              []Tag            `json:"tags" gorm:"many2many:applicationTags"`
	BusinessServiceID uint             `json:"-"`
	BusinessService   *BusinessService `json:"businessService"`
}

type Review struct {
	Model
	BusinessCriticality uint         `json:"businessCriticality" gorm:"notnull"`
	EffortEstimate      string       `json:"effortEstimate" gorm:"notnull"`
	ProposedAction      string       `json:"proposedAction" gorm:"notnull"`
	WorkPriority        uint         `json:"workPriority" gorm:"notnull"`
	Comments            string       `json:"comments"`
	Application         *Application `json:"application"`
	ApplicationID       uint         `json:"-"`
}

type Artifact struct {
	Model
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Location      string `json:"location"`
	ApplicationID uint   `json:"application"`
}
