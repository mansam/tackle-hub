package model

type BusinessService struct {
	Model
	Name          string        `json:"name" gorm:"notnull,unique"`
	Description   string        `json:"description"`
	Applications  []Application `json:"applications"`
	StakeholderID uint          `json:"owner"`
}

type StakeholderGroup struct {
	Model
	Name         string        `json:"name" gorm:"unique,index"`
	Description  string        `json:"description"`
	Stakeholders []Stakeholder `json:"stakeholders" gorm:"many2many:StakeholderGroups"`
}

type Stakeholder struct {
	Model
	Email            string             `json:"email" gorm:"notnull"`
	Name             string             `json:"displayName" gorm:"notnull"`
	Groups           []StakeholderGroup `json:"groups" gorm:"many2many:groups"`
	BusinessServices []BusinessService  `json:"businessServices"`
	JobFunctionID    uint               `json:"jobFunction" gorm:"notnull"`
}

type JobFunction struct {
	Model
	Role         string        `json:"role" gorm:"notnull,unique"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

type Tag struct {
	Model
	Name      string  `json:"name" gorm:"notnull"`
	Username  string  `json:"username"`
	TagType   TagType `json:"tagType" gorm:"notnull"`
	TagTypeID uint    `json:"-"`
}

type TagType struct {
	Model
	Name     string `json:"name" gorm:"notnull"`
	Username string `json:"username"`
	Rank     uint   `json:"rank"`
	Color    string `json:"colour"`
	Tags     []Tag  `json:"tags"`
}
