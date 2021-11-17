package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/models"
	"net/http"
)

//
// Routes
const (
	ApplicationsRoot = InventoryRoot + "/application"
	ApplicationRoot  = ApplicationsRoot + "/:" + ID
)

type ApplicationHandler struct {
	BaseHandler
}

func (h ApplicationHandler) AddRoutes(e *gin.Engine) {
	e.GET(ApplicationsRoot, h.List)
	e.GET(ApplicationsRoot+"/", h.List)
	e.POST(ApplicationsRoot, h.Create)
	e.GET(ApplicationRoot, h.Get)
	e.PUT(ApplicationRoot, h.Update)
	e.DELETE(ApplicationRoot, h.Delete)
}

// Get godoc
// @summary Get an application by ID.
// @description Get an application by ID.
// @tags get
// @produce json
// @success 200 {object} models.Application
// @router /application-inventory/application/:id [get]
// @param id path string true "Application ID"
func (h ApplicationHandler) Get(ctx *gin.Context) {
	model := models.Application{}
	id := ctx.Param(ID)
	result := h.DB.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all applications.
// @description List all applications.
// @tags list
// @produce json
// @success 200 {object} []models.Application
// @router /application-inventory/application [get]
func (h ApplicationHandler) List(ctx *gin.Context) {
	var list []models.Application
	page := NewPagination(ctx)
	result := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort).Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create an application.
// @description Create an application.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.Application
// @router /application-inventory/application [post]
// @param application body models.Application true "Application data"
func (h ApplicationHandler) Create(ctx *gin.Context) {
	model := models.Application{}
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
// @summary Delete an application.
// @description Delete an application.
// @tags delete
// @success 200 {object} models.Application
// @router /application-inventory/application/:id [delete]
// @param id path string true "Application id"
func (h ApplicationHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&models.Application{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update an application.
// @description Update an application.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.Application
// @router /application-inventory/application/:id [put]
// @param id path string true "Application id"
// @param application body models.Application true "Application data"
func (h ApplicationHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := models.Application{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&models.Application{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
