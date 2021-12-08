package importer

import (
	"context"
	"github.com/konveyor/tackle-hub/model"
	"gorm.io/gorm"
	"time"
)

//
//
type Manager struct {
	// DB
	DB *gorm.DB
}

//
//
func (m *Manager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <- ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				_ = m.importPending()
			}
		}
	}()
}

//
//
func (m *Manager) importPending() (err error) {
	list := []model.ApplicationImport{}
	result := m.DB.Find(&list, "processed = ?", false)
	if result.Error != nil {
		err = result.Error
		return
	}
	for _, imp := range list {
		imp.Processed = true
		_, ok := m.application(&imp)
		imp.IsValid = ok
		_ = m.DB.Save(&imp)
	}
	return
}

func (m *Manager) application(imp *model.ApplicationImport) (app *model.Application, ok bool) {
	app = &model.Application{}
	businessService := &model.BusinessService{}
	result := m.DB.Select("id").Where("name LIKE ?", imp.BusinessService).First(businessService)
	if result.Error != nil {
		return
	}
	app.BusinessService = businessService
	app.Name = imp.ApplicationName
	app.Description = imp.Description
	app.Comments = imp.Comments

	tags := []model.Tag{}
	db := m.DB.Preload("TagType")
	db.Find(&tags)

	for _, impTag := range imp.ImportTags {
		for _, tag := range tags {
			if tag.Name == impTag.Name && tag.TagType.Name == impTag.TagType {
				app.Tags = append(app.Tags, tag)
				continue
			}
		}
	}

	result = m.DB.Create(app)
	if result.Error != nil {
		return
	}
	ok = true
	return
}