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

type Artifact struct {
	Model
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Location      string `json:"location"`
	ApplicationID uint   `json:"application"`
}

type Dependency struct {
	Model
	ToID   uint         `json:"to"`
	To     *Application `json:"-" gorm:"foreignKey:to_id;constraint:OnDelete:CASCADE"`
	FromID uint         `json:"from"`
	From   *Application `json:"-" gorm:"foreignKey:from_id;constraint:OnDelete:CASCADE"`
}

type Review struct {
	Model
	BusinessCriticality uint         `json:"businessCriticality" gorm:"not null"`
	EffortEstimate      string       `json:"effortEstimate" gorm:"not null"`
	ProposedAction      string       `json:"proposedAction" gorm:"not null"`
	WorkPriority        uint         `json:"workPriority" gorm:"not null"`
	Comments            string       `json:"comments"`
	Application         *Application `json:"application"`
	ApplicationID       uint         `json:"-"`
}
