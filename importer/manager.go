package importer

import (
	"context"
	"fmt"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
	"gorm.io/gorm"
	"strings"
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
		var ok bool
		switch imp.RecordType1 {
		case api.RecordTypeApplication:
			ok = m.application(&imp)
		case api.RecordTypeDependency:
			ok = m.dependency(&imp)
		}
		imp.IsValid = ok
		imp.Processed = true
		_ = m.DB.Save(&imp)
	}
	return
}

func (m *Manager) dependency(imp *model.ApplicationImport) (ok bool) {
	app := &model.Application{}
	result := m.DB.Select("id").Where("name LIKE ?", imp.ApplicationName).First(app)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("Application '%s' does not exist.", imp.ApplicationName)
		return
	}

	dep := &model.Application{}
	result = m.DB.Select("id").Where("name LIKE ?", imp.Dependency).First(dep)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("Application dependency '%s' does not exist.", imp.Dependency)
		return
	}

	dependency := &model.Dependency{}
	switch strings.ToLower(imp.DependencyDirection) {
	case "northbound":
		dependency.FromID = dep.ID
		dependency.ToID = app.ID
	case "southbound":
		dependency.FromID = app.ID
		dependency.ToID = dep.ID
	}

	result = m.DB.Create(dependency)
	if result.Error != nil {
		imp.ErrorMessage = result.Error.Error()
		return
	}

	ok = true
	return
}

func (m *Manager) application(imp *model.ApplicationImport) (ok bool) {
	app := &model.Application{}
	businessService := &model.BusinessService{}
	result := m.DB.Select("id").Where("name LIKE ?", imp.BusinessService).First(businessService)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("BusinessService '%s' does not exist.", imp.BusinessService)
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
		imp.ErrorMessage = result.Error.Error()
		return
	}
	ok = true
	return
}