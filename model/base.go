package model

import (
	"bufio"
	"encoding/json"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/gorm"
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
		modelType := reflect.TypeOf(model)
		segments := strings.Split(modelType.String(), ".")
		fileName := strings.ToLower(segments[len(segments)-1]) + ".json"
		seedPath := path.Join(settings.Settings.DB.SeedPath, fileName)
		file, err := os.Open(seedPath)
		if err != nil {
			log.Info("No seed file found for type.", "type", modelType.String())
			continue
		}
		defer file.Close()
		created := 0
		failed := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			m := make(map[string]interface{})
			err = json.Unmarshal(scanner.Bytes(), &m)
			if err != nil {
				log.Info("Could not unmarshal record.", "type", modelType.String())
				failed++
				continue
			}
			result := db.Model(&model).Create(m)
			if result.Error != nil {
				log.Info("Could not create row.", "type", modelType.String(), "error", result.Error)
				continue
			}
			created++
		}
		log.Info("Complete.", "type", modelType.String(), "created", created, "failed", failed)
	}

	seeded := Seeded{}
	result = db.Create(&seeded)
	if result.Error != nil {
		log.Info("Failed to create seed record.")
	}
	log.Info("Database seeded.")
}
