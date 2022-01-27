package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

//
// BaseHandler base handler.
type BaseHandler struct {
	// DB
	DB *gorm.DB
	// k8s Client
	Client client.Client
}

// With database and k8s client.
func (h *BaseHandler) With(db *gorm.DB, client client.Client) {
	h.DB = db
	h.Client = client
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
	if errors.Is(err, os.ErrNotExist) {
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
	status := http.StatusInternalServerError
	sqliteErr := &sqlite3.Error{}

	if errors.As(err, sqliteErr) {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			status = http.StatusConflict
		}
	}

	ctx.JSON(
		status,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"Create failed.",
		"url",
		url)
}

//
// updateFailed handles Update() errors.
func (h *BaseHandler) updateFailed(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	sqliteErr := &sqlite3.Error{}

	if errors.As(err, sqliteErr) {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			status = http.StatusConflict
		}
	}

	ctx.JSON(
		status,
		gin.H{
			"error": err.Error(),
		})

	url := ctx.Request.URL.String()
	log.Error(
		err,
		"Update failed.",
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
		"Delete failed.",
		"url",
		url)
}

//
// bindFailed handles errors from BindJSON().
func (h *BaseHandler) bindFailed(ctx *gin.Context, err error) {
	ctx.JSON(
		http.StatusBadRequest,
		gin.H{
			"error": err.Error(),
		})
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
// listResponse selectively returns hal+json or plain json based on the "accept" header
func (h *BaseHandler) listResponse(ctx *gin.Context, kind string, resources interface{}, count int) {
	for _, accept := range ctx.Request.Header.Values("Accept") {
		if strings.Contains(accept, "application/hal+json") {
			ctx.Writer.Header().Set("Content-Type", "application/hal+json; charset=utf-8")
			hal := Hal{}
			hal.With(kind, resources, count)
			ctx.JSON(http.StatusOK, hal)
			return
		}
	}
	ctx.JSON(http.StatusOK, resources)
}

//
// Hal REST resource.
type Hal struct {
	Embedded   map[string]interface{} `json:"_embedded"`
	TotalCount int                    `json:"total_count"`
}

//
// With sets the embedded resource and count.
func (r *Hal) With(kind string, resources interface{}, total int) {
	r.Embedded = make(map[string]interface{})
	r.Embedded[kind] = resources
	r.TotalCount = total
}

//
// REST resource.
type Resource struct {
	ID         uint      `json:"id"`
	CreateUser string    `json:"createUser"`
	UpdateUser string    `json:"updateUser"`
	CreateTime time.Time `json:"createTime"`
}

//
// Update the resource with the model.
func (r *Resource) With(m *model.Model) {
	r.ID = m.ID
	r.CreateUser = m.CreateUser
	r.UpdateUser = m.UpdateUser
	r.CreateTime = m.CreateTime
}
