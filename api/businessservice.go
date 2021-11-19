package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Routes
const (
	BusinessServicesRoot = ControlsRoot + "/business-service"
	BusinessServiceRoot  = BusinessServicesRoot + "/:" + ID
)

//
// BusinessServiceHandler handles business-service routes.
type BusinessServiceHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h BusinessServiceHandler) AddRoutes(e *gin.Engine) {
	e.GET(BusinessServicesRoot, h.List)
	e.GET(BusinessServicesRoot+"/", h.List)
	e.POST(BusinessServicesRoot, h.Create)
	e.GET(BusinessServiceRoot, h.Get)
	e.PUT(BusinessServiceRoot, h.Update)
	e.DELETE(BusinessServiceRoot, h.Delete)
}

// Get godoc
// @summary Get a business service by ID.
// @description Get a business service by ID.
// @tags get
// @produce json
// @success 200 {object} model.BusinessService
// @router /controls/business-service/:id [get]
// @param id path string true "Business Service ID"
func (h BusinessServiceHandler) Get(ctx *gin.Context) {
	model := model.BusinessService{}
	id := ctx.Param(ID)
	db := h.DB.Preload("Applications")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all business services.
// @description List all business services.
// @tags list
// @produce json
// @success 200 {object} model.BusinessService
// @router /controls/business-service [get]
func (h BusinessServiceHandler) List(ctx *gin.Context) {
	var list []model.BusinessService
	page := NewPagination(ctx)
	db := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort)
	db = db.Preload("Applications")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a business service.
// @description Create a business service.
// @tags create
// @accept json
// @produce json
// @success 200 {object} model.BusinessService
// @router /controls/business-service [post]
// @param business_service body model.BusinessService true "Business service data"
func (h BusinessServiceHandler) Create(ctx *gin.Context) {
	model := model.BusinessService{}
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
// @summary Delete a business service.
// @description Delete a business service.
// @tags delete
// @success 200 {object} model.BusinessService
// @router /controls/business-service/:id [delete]
// @param id path string true "Business service ID"
func (h BusinessServiceHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.BusinessService{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a business service.
// @description Update a business service.
// @tags update
// @accept json
// @produce json
// @success 200 {object} model.BusinessService
// @router /controls/business-service/:id [put]
// @param id path string true "Business service ID"
// @param business_service body model.BusinessService true "Business service data"
func (h BusinessServiceHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := model.BusinessService{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&model.BusinessService{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
