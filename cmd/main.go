package main

import (
	"context"
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

	model.Seed(db,
		model.JobFunction{},
		model.TagType{},
		model.Tag{},
		model.StakeholderGroup{},
		model.Stakeholder{},
		model.BusinessService{},
		model.Application{},
		model.Review{},
	)

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
	handlerList := []api.Handler{
		&api.ApplicationHandler{},
		&api.BucketHandler{},
		&api.BusinessServiceHandler{},
		&api.DependencyHandler{},
		&api.ImportHandler{},
		&api.JobFunctionHandler{},
		&api.RepositoryHandler{},
		&api.IdentityHandler{},
		&api.ReviewHandler{},
		&api.StakeholderHandler{},
		&api.StakeholderGroupHandler{},
		&api.TagHandler{},
		&api.TagTypeHandler{},
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
	importManager := importer.Manager{
		DB: db,
	}
	importManager.Run(context.Background())
	err = router.Run()
}
