package model

type BusinessService struct {
	Model
	Name        string       `json:"name" gorm:"index;unique;not null"`
	Description string       `json:"description"`
	OwnerID     *uint        `json:"-" gorm:"index"`
	Owner       *Stakeholder `json:"owner"`
}

type StakeholderGroup struct {
	Model
	Name         string        `json:"name" gorm:"index;unique;not null"`
	Username     string        `json:"username"`
	Description  string        `json:"description"`
	Stakeholders []Stakeholder `json:"stakeholders" gorm:"many2many:sgStakeholder"`
}

type Stakeholder struct {
	Model
	DisplayName      string             `json:"displayName" gorm:"not null;"`
	Email            string             `json:"email" gorm:"index;unique;not null"`
	Groups           []StakeholderGroup `json:"stakeholderGroups" gorm:"many2many:sgStakeholder"`
	BusinessServices []BusinessService  `json:"businessServices" gorm:"foreignKey:OwnerID"`
	JobFunctionID    *uint              `json:"-" gorm:"index"`
	JobFunction      *JobFunction       `json:"jobFunction"`
}

type JobFunction struct {
	Model
	Username     string        `json:"username"`
	Role         string        `json:"role" gorm:"index;unique;not null"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

type Tag struct {
	Model
	Name      string  `json:"name" gorm:"uniqueIndex:tag_a;not null"`
	Username  string  `json:"username"`
	TagTypeID uint    `json:"-" gorm:"uniqueIndex:tag_a;index;not null"`
	TagType   TagType `json:"tagType"`
}

type TagType struct {
	Model
	Name     string `json:"name" gorm:"index;unique;not null"`
	Username string `json:"username"`
	Rank     uint   `json:"rank"`
	Color    string `json:"colour"`
	Tags     []Tag  `json:"tags"`
}
