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
// Manager for processing application imports.
type Manager struct {
	// DB
	DB *gorm.DB
}

//
// Run the manager.
func (m *Manager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				_ = m.processImports()
			}
		}
	}()
}

//
// processImports creates applications and dependencies from
// unprocessed imports.
func (m *Manager) processImports() (err error) {
	list := []model.Import{}
	db := m.DB.Preload("ImportTags")
	result := db.Find(&list, "processed = ?", false)
	if result.Error != nil {
		err = result.Error
		return
	}
	for _, imp := range list {
		var ok bool
		switch imp.RecordType1 {
		case api.RecordTypeApplication:
			ok = m.createApplication(&imp)
		case api.RecordTypeDependency:
			ok = m.createDependency(&imp)
		}
		imp.IsValid = ok
		imp.Processed = true
		result = m.DB.Save(&imp)
		if result.Error != nil {
			err = result.Error
			return
		}
	}
	return
}

//
// createDependency creates an application dependency from
// a dependency import record.
func (m *Manager) createDependency(imp *model.Import) (ok bool) {
	app := &model.Application{}
	result := m.DB.Select("id").Where("name LIKE ?", imp.ApplicationName).First(app)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("Application '%s' could not be found.", imp.ApplicationName)
		return
	}

	dep := &model.Application{}
	result = m.DB.Select("id").Where("name LIKE ?", imp.Dependency).First(dep)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("Application dependency '%s' could not be found.", imp.Dependency)
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

//
// createApplication creates an application from an
// application import record.
func (m *Manager) createApplication(imp *model.Import) (ok bool) {
	app := &model.Application{}
	businessService := &model.BusinessService{}
	result := m.DB.Select("id").Where("name LIKE ?", imp.BusinessService).First(businessService)
	if result.Error != nil {
		imp.ErrorMessage = fmt.Sprintf("BusinessService '%s' could not be found.", imp.BusinessService)
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
