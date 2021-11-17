package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/models"
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
	model := models.BinaryRepo{}
	id := ctx.Param(ID)
	result := h.DB.First(&model, "id = ?", id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all binary repositories.
// @description List all binary repositories.
// @tags get
// @produce json
// @success 200 {object} models.BinaryRepository
// @router /application-inventory/binary-repository [get]
func (h BinaryRepoHandler) List(ctx *gin.Context) {
	var list []models.BinaryRepo
	page := NewPagination(ctx)
	result := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort).Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
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
	model := models.BinaryRepo{}
	err := ctx.BindJSON(&model)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(&model)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
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
		h.deleteFailed(ctx, result.Error)
		return
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
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&models.BinaryRepo{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
