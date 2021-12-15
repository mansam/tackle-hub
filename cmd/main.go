package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/k8s"
	crd "github.com/konveyor/tackle-hub/k8s/api"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/settings"
	"github.com/konveyor/tackle-hub/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes/scheme"
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
		&model.Review{},
		&model.BusinessService{},
		&model.Dependency{},
		&model.JobFunction{},
		&model.Review{},
		&model.Stakeholder{},
		&model.StakeholderGroup{},
		&model.Tag{},
		&model.TagType{},
		&model.TaskReport{},
		&model.Task{})
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
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	handlerList := []api.Handler{
		&api.ApplicationHandler{},
		&api.BucketHandler{},
		&api.ReviewHandler{},
		&api.BusinessServiceHandler{},
		&api.DependencyHandler{},
		&api.JobFunctionHandler{},
		&api.ReviewHandler{},
		&api.TagHandler{},
		&api.TagTypeHandler{},
		&api.StakeholderHandler{},
		&api.StakeholderGroupHandler{},
		&api.TaskHandler{
			Client: client,
		},
		&api.AddonHandler{
			Client: client,
		},
	}
	for _, h := range handlerList {
		h.With(db)
		h.AddRoutes(router)
	}
	taskManager := task.Manager{
		Client: client,
		DB:     db,
	}
	taskManager.Run(context.Background())
	err = router.Run()
}
