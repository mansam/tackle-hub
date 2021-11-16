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
	BinaryReposRoot = InventoryRoot + "/binary-repository"
	BinaryRepoRoot  = BinaryReposRoot + "/:" + ID
)

type BinaryRepoHandler struct {
	BaseHandler
}

func (h BinaryRepoHandler) AddRoutes(e *gin.Engine) {
	e.GET(BinaryReposRoot, h.List)
	e.GET(BinaryReposRoot+"/", h.List)
	e.POST(BinaryReposRoot, h.Create)
	e.GET(BinaryRepoRoot, h.Get)
	e.PUT(BinaryRepoRoot, h.Update)
	e.DELETE(BinaryRepoRoot, h.Delete)
}

// Get godoc
// @summary Get a binary repository by ID.
// @description Get a binary repository by ID.
// @tags get
// @produce json
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository/:id [get]
// @param id path string true "Binary Repository ID"
func (h BinaryRepoHandler) Get(ctx *gin.Context) {
	binaryRepo := models.BinaryRepo{}
	id := ctx.Param(ID)
	result := h.DB.First(&binaryRepo, "id = ?", id)
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
	ctx.JSON(http.StatusOK, binaryRepo)
}

// List godoc
// @summary List all binary repositories.
// @description List all binary repositories.
// @tags get
// @produce json
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository [get]
func (h BinaryRepoHandler) List(ctx *gin.Context) {
	var binaryRepos []models.BinaryRepo
	result := h.DB.Find(&binaryRepos)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, binaryRepos)
}

// Create godoc
// @summary Create a binary repository.
// @description Create a binary repository.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository [post]
// @param binary_repository body models.BinaryRepository true "Binary Repository data"
func (h BinaryRepoHandler) Create(ctx *gin.Context) {
	binaryRepo := models.BinaryRepo{}
	err := ctx.BindJSON(&binaryRepo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Create(&binaryRepo)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, binaryRepo)
}

// Delete godoc
// @summary Delete a binary repository.
// @description Delete a binary repository.
// @tags delete
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository/:id [delete]
// @param id path string true "Binary Repository ID"
func (h BinaryRepoHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)

	result := h.DB.Delete(&models.BinaryRepo{}, "id = ?", id)
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
// @summary Update a binary repository.
// @description Update a binary repository.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository/:id [put]
// @param id path string true "Binary Repository ID"
// @param binary_repository body models.BinaryRepository true "Binary Repository data"
func (h BinaryRepoHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)

	updates := models.BinaryRepo{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.BinaryRepo{}).Where("id = ?", id).Omit("id").Updates(updates)
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
