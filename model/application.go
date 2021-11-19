package model

type Application struct {
	Model
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	BusinessServiceID uint    `json:"businessService"`
	Review            *Review `json:"review"`
	Comments          string  `json:"comments"`
	Tags              []Tag   `json:"tags" gorm:"many2many:ApplicationTags"`
}

type Review struct {
	Model
	Comments            string `json:"comments"`
	BusinessCriticality uint   `json:"businessCriticality" gorm:"notnull"`
	EffortEstimate      string `json:"effortEstimate" gorm:"notnull"`
	ProposedAction      string `json:"proposedAction" gorm:"notnull"`
	WorkPriority        uint   `json:"workPriority" gorm:"notnull"`
	ApplicationID       uint   `json:"application"`
}
