package model

import "fmt"

type Application struct {
	Model
	Name              string           `json:"name" gorm:"not null"`
	Description       string           `json:"description"`
	Review            *Review          `json:"review"`
	Comments          string           `json:"comments"`
	Tags              []Tag            `json:"tags" gorm:"many2many:applicationTags"`
	Identities        []Identity       `json:"identities" gorm:"many2many:appIdentity"`
	BusinessServiceID uint             `json:"-" gorm:"index"`
	BusinessService   *BusinessService `json:"businessService"`
}

type Dependency struct {
	Model
	ToID   uint         `json:"to" gorm:"index"`
	To     *Application `json:"-" gorm:"foreignKey:to_id;constraint:OnDelete:CASCADE"`
	FromID uint         `json:"from" gorm:"index"`
	From   *Application `json:"-" gorm:"foreignKey:from_id;constraint:OnDelete:CASCADE"`
}

type Repository struct {
	Model
	Kind          string `json:"kind"`
	URL           string `json:"url"`
	Branch        string `json:"branch"`
	Tag           string `json:"tag"`
	Path          string `json:"path" gorm:"default:/"`
	ApplicationID uint   `json:"application" gorm:"index;unique"`
}

type Review struct {
	Model
	BusinessCriticality uint         `json:"businessCriticality" gorm:"not null"`
	EffortEstimate      string       `json:"effortEstimate" gorm:"not null"`
	ProposedAction      string       `json:"proposedAction" gorm:"not null"`
	WorkPriority        uint         `json:"workPriority" gorm:"not null"`
	Comments            string       `json:"comments"`
	Application         *Application `json:"application"`
	ApplicationID       uint         `json:"-" gorm:"index"`
}

type Import struct {
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
	ImportSummary       ImportSummary `json:"-"`
	ImportSummaryID     uint          `json:"-" gorm:"index"`
	Processed           bool          `json:"-"`
	ImportTags          []ImportTag   `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (r *Import) AsMap() (m map[string]interface{}) {
	m = make(map[string]interface{})
	m["filename"] = r.Filename
	m["applicationName"] = r.ApplicationName
	// "Application Name" is necessary in order for
	// the UI to display the error report correctly.
	m["Application Name"] = r.ApplicationName
	m["businessService"] = r.BusinessService
	m["comments"] = r.Comments
	m["dependency"] = r.Dependency
	m["dependencyDirection"] = r.DependencyDirection
	m["description"] = r.Description
	m["errorMessage"] = r.ErrorMessage
	m["isValid"] = r.IsValid
	m["processed"] = r.Processed
	m["recordType1"] = r.RecordType1
	for i, tag := range r.ImportTags {
		m[fmt.Sprintf("tagType%v", i+1)] = tag.TagType
		m[fmt.Sprintf("tag%v", i+1)] = tag.Name
	}
	return
}

type ImportSummary struct {
	Model
	Content      []byte   `json:"-"`
	Filename     string   `json:"filename"`
	ImportStatus string   `json:"importStatus" gorm:"column:importStatus"`
	Imports      []Import `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

type ImportTag struct {
	Model
	Name     string
	TagType  string
	ImportID uint `gorm:"index"`
	Import   *Import
}
