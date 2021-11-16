package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Resource struct {
	ID        string `sql:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *Resource) BeforeCreate(_ *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}

type Application struct {
	Resource
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Comments          string          `json:"comments"`
	DependsOn         []*Application  `json:"depends_on" gorm:"many2many:application_dependencies"`
	Tags              []Tag           `json:"tags" gorm:"many2many:application_tags"`
	BusinessServiceID string          `json:"business_service_id" gorm:"notnull" binding:"required"`
	BusinessService   BusinessService `json:"businessService"`
}

type BinaryRepo struct {
	Resource
	Type          string `json:"name" gorm:"notnull" binding:"required" validate:"oneof=mvn"`
	URL           string `json:"url" gorm:"notnull" binding:"required"`
	Group         string `json:"group" gorm:"notnull" binding:"required"`
	Artifact      string `json:"artifact" gorm:"notnull" binding:"required"`
	Version       string `json:"version" gorm:"notnull" binding:"required"`
	ApplicationID string `json:"application_id"`
	Application   Application
}

type BusinessService struct {
	Resource
	Name        string      `json:"name" gorm:"notnull,unique" validate:"required"`
	Description string      `json:"description"`
	OwnerID     string      `json:"owner_id"`
	Owner       Stakeholder `json:"owner"`
}

type Group struct {
	Resource
	Name         string         `json:"name" gorm:"unique,index" validate:"required,alphanum,min=6,max=32"`
	Description  string         `json:"description"`
	Stakeholders []*Stakeholder `json:"stakeholders" gorm:"many2many:stakeholder_groups"`
}

type JobFunction struct {
	Resource
	Role         string        `json:"role" gorm:"notnull,unique"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

type Review struct {
	Resource
	Comments            string `json:"comments"`
	BusinessCriticality uint   `json:"businessCriticality" gorm:"notnull" binding:"required"`
	EffortEstimate      string `json:"effortEstimate" gorm:"notnull" binding:"required"`
	ProposedAction      string `json:"proposedAction" gorm:"notnull" binding:"required"`
	WorkPriority        uint   `json:"workPriority" gorm:"notnull" binding:"required"`
	ApplicationID       string `json:"application_id"`
	Application         Application
}

type SourceRepo struct {
	Resource
	Type          string `json:"name" gorm:"notnull" binding:"required" validate:"oneof=git svn"`
	URL           string `json:"url" gorm:"notnull" binding:"required"`
	Branch        string `json:"branch" gorm:"notnull" binding:"required"`
	ApplicationID string `json:"application_id"`
	Application   Application
}

type Stakeholder struct {
	Resource
	Email         string   `json:"email" gorm:"notnull" binding:"required,email"`
	DisplayName   string   `json:"displayName" gorm:"notnull" binding:"required"`
	Groups        []*Group `json:"groups" gorm:"many2many:stakeholder_groups"`
	JobFunctionID string   `json:"job_function_id" gorm:"notnull" binding:"required"`
	JobFunction   JobFunction
}

type Tag struct {
	Resource
	Name      string `json:"name" gorm:"notnull" binding:"required"`
	TagTypeID string `json:"tag_type_id" gorm:"notnull" binding:"required"`
	TagType   TagType
}

type TagType struct {
	Resource
	Name   string `json:"name" gorm:"notnull" binding:"required"`
	Rank   uint   `json:"rank"`
	Colour string `json:"colour"`
	Tags   []Tag  `json:"tags"`
}
