package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/models"
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

func dbPath() string {
	base, ok := os.LookupEnv(DatabasePathEnv)
	if !ok {
		log.Fatal(fmt.Sprintf("%s not set, aborting.", DatabasePathEnv))
	}
	return path.Join(base, DatabaseFileName)
}

func Setup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(
		&models.Application{},
		&models.BinaryRepo{},
		&models.BusinessService{},
		&models.Group{},
		&models.JobFunction{},
		&models.Review{},
		&models.SourceRepo{},
		&models.Tag{},
		&models.TagType{},
		&models.Stakeholder{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := Setup()
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	handlerList := []api.Handler{
		&api.ApplicationHandler{},
		&api.BinaryRepoHandler{},
		&api.BusinessServiceHandler{},
		&api.GroupHandler{},
		&api.JobFunctionHandler{},
		&api.ReviewHandler{},
		&api.SourceRepoHandler{},
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
