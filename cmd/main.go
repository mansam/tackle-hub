package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/k8s"
	crd "github.com/konveyor/tackle-hub/k8s/api"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"os"
)

const (
	DatabasePathEnv  = "DB_PATH"
)

//
// dbPath builds DB path.
func dbPath() (path string) {
	path, found := os.LookupEnv(DatabasePathEnv)
	if !found {
		path = "/tmp/tackle.db"
	}

	return
}

//
// Setup the DB and models.
func Setup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(
		&model.Application{},
		&model.Artifact{},
		&model.Review{},
		&model.BusinessService{},
		&model.StakeholderGroup{},
		&model.JobFunction{},
		&model.Tag{},
		&model.TagType{},
		&model.Stakeholder{},
		&model.TaskReport{},
		&model.Task{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//
// buildScheme adds CRDs to the k8s scheme.
func buildScheme() {
	err := crd.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatal(err, "Add CRD failed.")
		os.Exit(1)
	}
}

//
// main.
func main() {
	buildScheme()
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	db := Setup()
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	handlerList := []api.Handler{
		&api.ApplicationHandler{},
		&api.ArtifactHandler{},
		&api.ReviewHandler{},
		&api.BusinessServiceHandler{},
		&api.StakeholderGroupHandler{},
		&api.JobFunctionHandler{},
		&api.TagHandler{},
		&api.TagTypeHandler{},
		&api.StakeholderHandler{},
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
	if err != nil {
		log.Fatal(err)
	}
}
