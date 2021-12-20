package model

type BusinessService struct {
	Model
	Name        string       `json:"name" gorm:"not null; unique"`
	Description string       `json:"description"`
	OwnerID     *uint        `json:"-"`
	Owner       *Stakeholder `json:"owner"`
}

type StakeholderGroup struct {
	Model
	Name         string        `json:"name" gorm:"unique,index"`
	Username     string        `json:"username"`
	Description  string        `json:"description"`
	Stakeholders []Stakeholder `json:"stakeholders" gorm:"many2many:sgStakeholder"`
}

type Stakeholder struct {
	Model
	DisplayName      string             `json:"displayName" gorm:"column:displayName;not null"`
	Email            string             `json:"email" gorm:"not null"`
	Groups           []StakeholderGroup `json:"stakeholderGroups" gorm:"many2many:sgStakeholder"`
	BusinessServices []BusinessService  `json:"businessServices" gorm:"foreignKey:owner_id"`
	JobFunctionID    *uint              `json:"-"`
	JobFunction      *JobFunction       `json:"jobFunction"`
}

type JobFunction struct {
	Model
	Username     string        `json:"username"`
	Role         string        `json:"role" gorm:"not null; unique"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

type Tag struct {
	Model
	Name      string  `json:"name" gorm:"not null"`
	Username  string  `json:"username"`
	TagTypeID uint    `json:"-" gorm:"not null"`
	TagType   TagType `json:"tagType"`
}

type TagType struct {
	Model
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username"`
	Rank     uint   `json:"rank"`
	Color    string `json:"colour"`
	Tags     []Tag  `json:"tags"`
}
