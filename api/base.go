package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"gorm.io/gorm"
	"strconv"
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

// Pagination Defaults
const (
	Limit  = 20
	Offset = 0
	Sort   = "created_at asc"
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

type Pagination struct {
	Limit  int
	Offset int
	Sort   string
}

func NewPagination(ctx *gin.Context) Pagination {
	limit, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		limit = Limit
	}
	offset, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		offset = Offset
	}
	sort := ctx.Query("sort")
	if sort == "" {
		sort = Sort
	}
	return Pagination{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
}
