package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"gorm.io/gorm"
)

var log = logging.WithName("api")

// Error messages
const (
	MsgInternalServerError = "internal server error"
	MsgNotFound            = "not found"
	MsgBadRequest          = "bad request"
)

// Routes
const (
	InventoryRoot = "/application-inventory"
	ControlsRoot  = "/controls"
)

// Params
const (
	ID = "id"
)

type Handler interface {
	AddRoutes(e *gin.Engine)
	Get(ctx *gin.Context)
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type BaseHandler struct {
	DB *gorm.DB
}
