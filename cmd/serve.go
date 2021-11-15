package main

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/db"
	"github.com/konveyor/tackle-hub/handlers"
	"log"
)

func main() {
	db.Setup()
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	applicationHandler := handlers.ApplicationHandler{}
	applicationHandler.AddRoutes(r)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
