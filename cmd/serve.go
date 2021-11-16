package main

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/db"
	"github.com/konveyor/tackle-hub/handlers"
	"log"
)

func main() {
	db.Setup()
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	handlerList := []handlers.Handler{
		&handlers.ApplicationHandler{},
		&handlers.BinaryRepoHandler{},
		&handlers.BusinessServiceHandler{},
		&handlers.GroupHandler{},
		&handlers.JobFunctionHandler{},
		&handlers.ReviewHandler{},
		&handlers.RoleHandler{},
		&handlers.SourceRepoHandler{},
		&handlers.TagHandler{},
		&handlers.TagTypeHandler{},
		&handlers.StakeholderHandler{},
	}
	for _, h := range handlerList {
		h.AddRoutes(e)
	}

	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
