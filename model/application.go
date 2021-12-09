package model

import (
	"fmt"
)

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

type Artifact struct {
	Model
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Location      string `json:"location"`
	ApplicationID uint   `json:"application"`
}

type ImportSummary struct {
	Model
	Filename           string              `json:"filename"`
	ImportStatus       string              `json:"importStatus" gorm:"column:importStatus"`
	ApplicationImports []ApplicationImport `json:"-"`
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
	ImportSummary       ImportSummary `json:"-"`
	ImportSummaryID     uint          `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Processed           bool          `json:"-"`
	ImportTags          []ImportTag   `json:"-"`
}

func (r *ApplicationImport) AsMap() (m map[string]interface{}) {
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

type ImportTag struct {
	Model
	Name                string
	TagType             string
	Order               uint
	ApplicationImportID uint `gorm:"constraint:OnDelete:CASCADE"`
	ApplicationImport   *ApplicationImport
}
