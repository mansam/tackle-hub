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
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	handlerList := []api.Handler{
		&api.ApplicationHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.BinaryRepoHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.BusinessServiceHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.GroupHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.JobFunctionHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.ReviewHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.SourceRepoHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.TagHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.TagTypeHandler{BaseHandler: api.BaseHandler{DB: db}},
		&api.StakeholderHandler{BaseHandler: api.BaseHandler{DB: db}},
	}
	for _, h := range handlerList {
		h.AddRoutes(e)
	}

	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
