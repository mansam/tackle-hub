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
		_ = m.DB.Save(&imp)
	}
	return
}

