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
	BusinessServicesRoot = ControlsRoot + "/business-service"
	BusinessServiceRoot  = BusinessServicesRoot + "/:" + ID
)

type BusinessServiceHandler struct {
	BaseHandler
}

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
// @success 200 {object} models.BusinessService
// @router /controls/business-service/:id [get]
// @param id path string true "Business Service ID"
func (h BusinessServiceHandler) Get(ctx *gin.Context) {
	model := models.BusinessService{}
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
// @summary List all business services.
// @description List all business services.
// @tags list
// @produce json
// @success 200 {object} models.BusinessService
// @router /controls/business-service [get]
func (h BusinessServiceHandler) List(ctx *gin.Context) {
	var list []models.BusinessService
	page := NewPagination(ctx)
	result := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort).Find(&list)
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
// @summary Create a business service.
// @description Create a business service.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.BusinessService
// @router /controls/business-service [post]
// @param business_service body models.BusinessService true "Business service data"
func (h BusinessServiceHandler) Create(ctx *gin.Context) {
	model := models.BusinessService{}
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
// @summary Delete a business service.
// @description Delete a business service.
// @tags delete
// @success 200 {object} models.BusinessService
// @router /controls/business-service/:id [delete]
// @param id path string true "Business service ID"
func (h BusinessServiceHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&models.BusinessService{}, "id = ?", id)
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
// @summary Update a business service.
// @description Update a business service.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.BusinessService
// @router /controls/business-service/:id [put]
// @param id path string true "Business service ID"
// @param business_service body models.BusinessService true "Business service data"
func (h BusinessServiceHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := models.BusinessService{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.BusinessService{}).Where("id = ?", id).Omit("id").Updates(updates)
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
