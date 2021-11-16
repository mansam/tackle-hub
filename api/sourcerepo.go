package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/models"
	"gorm.io/gorm"
	"net/http"
)

//
// Routes
const (
	SourceReposRoot = InventoryRoot + "/source-repository"
	SourceRepoRoot  = SourceReposRoot + "/:" + ID
)

type SourceRepoHandler struct {
	BaseHandler
}

func (h SourceRepoHandler) AddRoutes(e *gin.Engine) {
	e.GET(SourceReposRoot, h.List)
	e.GET(SourceReposRoot+"/", h.List)
	e.POST(SourceReposRoot, h.Create)
	e.GET(SourceRepoRoot, h.Get)
	e.PUT(SourceRepoRoot, h.Update)
	e.DELETE(SourceRepoRoot, h.Delete)
}

// Get godoc
// @summary Get a source repository by ID.
// @description Get a source repository by ID.
// @tags get
// @produce json
// @success 200 {object} models.SourceRepository
// @router /application-inventory/source-repository/:id [get]
// @param id path string true "Source Repository ID"
func (h SourceRepoHandler) Get(ctx *gin.Context) {
	model := models.SourceRepo{}
	id := ctx.Param(ID)
	result := h.DB.First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": MsgNotFound,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all source repositories.
// @description List all source repositories.
// @tags get
// @produce json
// @success 200 {object} models.SourceRepository
// @router /application-inventory/source-repository [get]
func (h SourceRepoHandler) List(ctx *gin.Context) {
	var list []models.SourceRepo
	result := h.DB.Find(&list)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a source repository.
// @description Create a source repository.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.SourceRepository
// @router /application-inventory/source-repository [post]
// @param source_repository body models.SourceRepository true "Source Repository data"
func (h SourceRepoHandler) Create(ctx *gin.Context) {
	model := models.SourceRepo{}
	err := ctx.BindJSON(&model)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Create(&model)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, model)
}

// Delete godoc
// @summary Delete a source repository.
// @description Delete a source repository.
// @tags delete
// @success 200 {object} models.SourceRepository
// @router /application-inventory/source-repository/:id [delete]
// @param id path string true "Source Repository ID"
func (h SourceRepoHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)

	result := h.DB.Delete(&models.SourceRepo{}, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusOK)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a source repository.
// @description Update a source repository.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.SourceRepository
// @router /application-inventory/source-repository/:id [put]
// @param id path string true "Source Repository ID"
// @param source_repository body models.SourceRepository true "Source Repository data"
func (h SourceRepoHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)

	updates := models.SourceRepo{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.SourceRepo{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": MsgNotFound,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.Status(http.StatusOK)
}
