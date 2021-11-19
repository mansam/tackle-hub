package model

type BusinessService struct {
	Model
	Name          string        `json:"name" gorm:"notnull,unique"`
	Description   string        `json:"description"`
	StakeholderID uint          `json:"-"`
	Stakeholder   *Stakeholder  `json:"owner"`
}

type StakeholderGroup struct {
	Model
	Name         string        `json:"name" gorm:"unique,index"`
	Description  string        `json:"description"`
	Stakeholders []Stakeholder `json:"stakeholders" gorm:"many2many:StakeholderGroup"`
}

type Stakeholder struct {
	Model
	Name             string             `json:"displayName" gorm:"notnull"`
	Email            string             `json:"email" gorm:"notnull"`
	Groups           []StakeholderGroup `json:"groups" gorm:"many2many:StakeholderGroup"`
	BusinessServices []BusinessService  `json:"businessServices"`
	JobFunctionID    uint               `json:"-" gorm:"notnull"`
	JobFunction      *JobFunction       `json:"jobFunction"`
}

type JobFunction struct {
	Model
	username     string        `json:"username"`
	Role         string        `json:"role" gorm:"notnull,unique"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

type Tag struct {
	Model
	Name      string  `json:"name" gorm:"notnull"`
	Username  string  `json:"username"`
	TagTypeID uint    `json:"-" gorm:"notnull"`
	TagType   TagType `json:"tagType"`
}

type TagType struct {
	Model
	Name     string `json:"name" gorm:"notnull"`
	Username string `json:"username"`
	Rank     uint   `json:"rank"`
	Color    string `json:"colour"`
	Tags     []Tag  `json:"tags"`
}
