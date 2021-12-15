package model

import (
	"gorm.io/gorm"
	"os"
)

type Bucket struct {
	Model
	Name          string `json:"name"`
	Path          string `json:"path"`
	ApplicationID uint   `json:"application"`
}

func (m *Bucket) AfterDelete(db *gorm.DB) (err error) {
	err = os.RemoveAll(m.Path)
	return
}
