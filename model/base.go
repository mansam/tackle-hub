package model

import (
	"encoding/json"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

var log = logging.WithName("model")

//
// Model Base model.
type Model struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreateUser string    `json:"createUser"`
	UpdateUser string    `json:"updateUser"`
	CreateTime time.Time `json:"createTime" gorm:"column:createTime;autoCreateTime"`
}

//
// Seeded model.
type Seeded struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

//
// Seed the database with the contents of json
// files contained in DB_SEED_PATH.
func Seed(db *gorm.DB, models ...interface{}) {
	result := db.Find(&Seeded{})
	if result.RowsAffected != 0 {
		log.Info("Database already seeded, skipping.")
		return
	}

	for _, model := range models {
		kind := reflect.TypeOf(model).String()
		segments := strings.Split(kind, ".")
		fileName := strings.ToLower(segments[len(segments)-1]) + ".json"
		filePath := path.Join(settings.Settings.DB.SeedPath, fileName)
		file, err := os.Open(filePath)
		if err != nil {
			log.Info("No seed file found for type.", "type", kind)
			continue
		}
		defer file.Close()
		jsonBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Info("Could not read bytes from file.", "type", kind)
		}

		var m []map[string]interface{}
		err = json.Unmarshal(jsonBytes, &m)
		if err != nil {
			log.Info("Could not unmarshal records.", "type", kind)
			continue
		}
		created := 0
		failed := 0
		for i := range m {
			result := db.Model(&model).Create(m[i])
			if result.Error != nil {
				log.Info("Could not create row.", "type", kind, "error", result.Error)
				failed++
				continue
			}
			created++
		}
		log.Info("Complete.", "type", kind, "created", created, "failed", failed)
	}

	seeded := Seeded{}
	result = db.Create(&seeded)
	if result.Error != nil {
		log.Info("Failed to create seed record.")
	}
	log.Info("Database seeded.")
}
