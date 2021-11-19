package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
)

const (
	DatabasePathEnv  = "TACKLE_HUB_DB_PATH"
	DatabaseFileName = "tackle-hub.sqlite"
)

//
// dbPath builds DB path.
func dbPath() string {
	dir, found := os.LookupEnv(DatabasePathEnv)
	if !found {
		log.Fatal(fmt.Sprintf("%s not set, aborting.", DatabasePathEnv))
	}

	return path.Join(dir, DatabaseFileName)
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
		&model.Review{},
		&model.BusinessService{},
		&model.StakeholderGroup{},
		&model.JobFunction{},
		&model.Tag{},
		&model.TagType{},
		&model.Stakeholder{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//
// main.
func main() {
	db := Setup()
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	handlerList := []api.Handler{
		&api.ApplicationHandler{},
		&api.ReviewHandler{},
		&api.BusinessServiceHandler{},
		&api.StakeholderGroupHandler{},
		&api.JobFunctionHandler{},
		&api.TagHandler{},
		&api.TagTypeHandler{},
		&api.StakeholderHandler{},
	}
	for _, h := range handlerList {
		h.With(db)
		h.AddRoutes(router)
	}
	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
