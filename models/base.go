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
