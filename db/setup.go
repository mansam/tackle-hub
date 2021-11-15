package db

import (
	"fmt"
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

var DB *gorm.DB

func dbPath() string {
	base, ok := os.LookupEnv(DatabasePathEnv)
	if !ok {
		log.Fatal(fmt.Sprintf("%s not set, aborting.", DatabasePathEnv))
	}
	return path.Join(base, DatabaseFileName)
}

func Setup() {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&models.Application{})
	if err != nil {
		log.Fatal(err)
	}
}
