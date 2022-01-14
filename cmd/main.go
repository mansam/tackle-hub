package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/importer"
	"github.com/konveyor/tackle-hub/k8s"
	crd "github.com/konveyor/tackle-hub/k8s/api"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/settings"
	"github.com/konveyor/tackle-hub/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path"
	"reflect"
	"strings"
	"syscall"
)

//
// DB constants
const (
	ConnectionString = "file:%s?_foreign_keys=yes"
)

var Settings = &settings.Settings

var log = logging.WithName("hub")

func init() {
	_ = Settings.Load()
}

//
// Setup the DB and models.
func Setup() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(fmt.Sprintf(ConnectionString, Settings.DB.Path)), &gorm.Config{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(
		&model.Application{},
		&model.Bucket{},
		&model.BusinessService{},
		&model.Dependency{},
		&model.Import{},
		&model.ImportSummary{},
		&model.ImportTag{},
		&model.JobFunction{},
		&model.Repository{},
		&model.Identity{},
		&model.Review{},
		&model.Seeded{},
		&model.Stakeholder{},
		&model.StakeholderGroup{},
		&model.Tag{},
		&model.TagType{},
		&model.TaskReport{},
		&model.Task{})
	if err != nil {
		return
	}

	err = Seed(db, model.All())
	if err != nil {
		return
	}

	return
}

//
// buildScheme adds CRDs to the k8s scheme.
func buildScheme() (err error) {
	err = crd.AddToScheme(scheme.Scheme)
	return
}

//
// main.
func main() {
	log.Info("Started", "settings", Settings)
	var err error
	defer func() {
		if err != nil {
			log.Trace(err)
		}
	}()
	syscall.Umask(0)
	err = buildScheme()
	if err != nil {
		return
	}
	client, err := k8s.NewClient()
	if err != nil {
		return
	}
	db, err := Setup()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	for _, h := range api.All() {
		h.With(db, client)
		h.AddRoutes(router)
	}
	taskManager := task.Manager{
		Client: client,
		DB:     db,
	}
	taskManager.Run(context.Background())
	importManager := importer.Manager{
		DB: db,
	}
	importManager.Run(context.Background())
	err = router.Run()
}

//
// Seed the database with the contents of json
// files contained in DB_SEED_PATH.
func Seed(db *gorm.DB, models []interface{}) (err error) {
	result := db.Find(&model.Seeded{})
	if result.RowsAffected != 0 {
		log.Info("Database already seeded, skipping.")
		return
	}

	for _, m := range models {
		err = func() (err error) {
			kind := reflect.TypeOf(m).Name()
			fileName := strings.ToLower(kind) + ".json"
			filePath := path.Join(settings.Settings.DB.SeedPath, fileName)
			file, err := os.Open(filePath)
			if err != nil {
				log.Info("Could not open seed file.", "model", kind, "path", filePath)
				err = nil
				return
			}
			defer file.Close()
			jsonBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return
			}

			var unmarshalled []map[string]interface{}
			err = json.Unmarshal(jsonBytes, &unmarshalled)
			if err != nil {
				return
			}
			for i := range unmarshalled {
				result := db.Model(&m).Create(unmarshalled[i])
				if result.Error != nil {
					err = result.Error
					return
				}
			}
			return
		}()
		if err != nil {
			return
		}
	}

	seeded := model.Seeded{}
	result = db.Create(&seeded)
	if result.Error != nil {
		err = result.Error
		return
	}
	log.Info("Database seeded.")
	return
}
