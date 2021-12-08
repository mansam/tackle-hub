package model

type Application struct {
	Model
	Name              string           `json:"name" gorm:"not null"`
	Description       string           `json:"description"`
	Review            *Review          `json:"review"`
	Comments          string           `json:"comments"`
	Tags              []Tag            `json:"tags" gorm:"many2many:applicationTags"`
	BusinessServiceID uint             `json:"-"`
	BusinessService   *BusinessService `json:"businessService"`
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

type Artifact struct {
	Model
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Location      string `json:"location"`
	ApplicationID uint   `json:"application"`
}

type ImportSummary struct {
	Model
	Filename     string `json:"filename"`
	ValidCount   int    `json:"validCount"`
	InvalidCount int    `json:"invalidCount"`
	ImportStatus string `json:"importStatus"`
}

type ApplicationImport struct {
	Model
	Filename            string        `json:"filename"`
	ApplicationName     string        `json:"applicationName"`
	BusinessService     string        `json:"businessService"`
	Comments            string        `json:"comments"`
	Dependency          string        `json:"dependency"`
	DependencyDirection string        `json:"dependencyDirection"`
	Description         string        `json:"description"`
	ErrorMessage        string        `json:"errorMessage"`
	IsValid             bool          `json:"isValid"`
	RecordType1         string        `json:"recordType1"`
	Tag1                string        `json:"tag1"`
	Tag2                string        `json:"tag2"`
	Tag3                string        `json:"tag3"`
	Tag4                string        `json:"tag4"`
	Tag5                string        `json:"tag5"`
	Tag6                string        `json:"tag6"`
	Tag7                string        `json:"tag7"`
	Tag8                string        `json:"tag8"`
	Tag9                string        `json:"tag9"`
	Tag10               string        `json:"tag10"`
	Tag11               string        `json:"tag11"`
	Tag12               string        `json:"tag12"`
	Tag13               string        `json:"tag13"`
	Tag14               string        `json:"tag14"`
	Tag15               string        `json:"tag15"`
	Tag16               string        `json:"tag16"`
	Tag17               string        `json:"tag17"`
	Tag18               string        `json:"tag18"`
	Tag19               string        `json:"tag19"`
	Tag20               string        `json:"tag20"`
	TagType1            string        `json:"tagType1"`
	TagType2            string        `json:"tagType2"`
	TagType3            string        `json:"tagType3"`
	TagType4            string        `json:"tagType4"`
	TagType5            string        `json:"tagType5"`
	TagType6            string        `json:"tagType6"`
	TagType7            string        `json:"tagType7"`
	TagType8            string        `json:"tagType8"`
	TagType9            string        `json:"tagType9"`
	TagType10           string        `json:"tagType10"`
	TagType11           string        `json:"tagType11"`
	TagType12           string        `json:"tagType12"`
	TagType13           string        `json:"tagType13"`
	TagType14           string        `json:"tagType14"`
	TagType15           string        `json:"tagType15"`
	TagType16           string        `json:"tagType16"`
	TagType17           string        `json:"tagType17"`
	TagType18           string        `json:"tagType18"`
	TagType19           string        `json:"tagType19"`
	TagType20           string        `json:"tagType20"`
	ImportSummary       ImportSummary `json:"-"`
	ImportSummaryID     uint          `json:"-"`
	Processed           bool          `json:"-"`
	Tags              []Tag            `json:"tags" gorm:"many2many:applicationTags"`
}

func (r *ApplicationImport) AsApplication() (app Application) {
	return
}
