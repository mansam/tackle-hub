package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var log = logging.WithName("api")

//
// Routes
const (
	InventoryRoot = "/application-inventory"
	ControlsRoot  = "/controls"
)

//
// Params
const (
	ID   = "id"
	Name = "name"
)

//
// Pagination Defaults
const (
	Limit  = 20
	Offset = 0
	Sort   = "id asc"
)

//
// Handler.
type Handler interface {
	With(*gorm.DB)
	AddRoutes(e *gin.Engine)
}

//
// BaseHandler base handler.
type BaseHandler struct {
	// DB
	DB *gorm.DB
}

// With database.
func (h *BaseHandler) With(db *gorm.DB) {
	h.DB = db
}

//
// getFailed handles Get() errors.
func (h *BaseHandler) getFailed(ctx *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			})
		return
	}
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"Get failed.",
		"url",
		url)
}

//
// listFailed handles List() errors.
func (h *BaseHandler) listFailed(ctx *gin.Context, err error) {
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"List failed.",
		"url",
		url)
}

//
// createFailed handles Create() errors.
func (h *BaseHandler) createFailed(ctx *gin.Context, err error) {
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"List failed.",
		"url",
		url)
}

//
// updateFailed handles Update() errors.
func (h *BaseHandler) updateFailed(ctx *gin.Context, err error) {
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"List failed.",
		"url",
		url)
}

//
// deleteFailed handles Delete() errors.
func (h *BaseHandler) deleteFailed(ctx *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusOK)
		return
	}
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"Get failed.",
		"url",
		url)
}

//
// preLoad update DB to pre-load fields.
func (h *BaseHandler) preLoad(db *gorm.DB, fields ...string) (tx *gorm.DB) {
	tx = db
	for _, f := range fields {
		tx = tx.Preload(f)
	}

	return
}

//
// Pagination provides pagination and sorting.
type Pagination struct {
	Limit  int
	Offset int
	Sort   string
}

//
// apply pagination.
func (p *Pagination) apply(db *gorm.DB) (tx *gorm.DB) {
	tx = db.Offset(p.Offset).Limit(p.Limit)
	tx = tx.Order(p.Sort)
	return
}

//
// NewPagination factory.
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
