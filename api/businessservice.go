package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	BusinessServiceKind = "business-service"
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
// @success 200 {object} api.BusinessService
// @router /controls/business-service/{id} [get]
// @param id path string true "Business Service ID"
func (h BusinessServiceHandler) Get(ctx *gin.Context) {
	m := &model.BusinessService{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Owner")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	resource := BusinessService{}
	resource.With(m)
	ctx.JSON(http.StatusOK, resource)
}

// List godoc
// @summary List all business services.
// @description List all business services.
// @tags list
// @produce json
// @success 200 {object} api.BusinessService
// @router /controls/business-service [get]
func (h BusinessServiceHandler) List(ctx *gin.Context) {
	var count int64
	var list []model.BusinessService
	h.DB.Model(&model.BusinessService{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Owner")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []BusinessService{}
	for i := range list {
		r := BusinessService{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, BusinessServiceKind, resources, int(count))
}

// Create godoc
// @summary Create a business service.
// @description Create a business service.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.BusinessService
// @router /controls/business-service [post]
// @param business_service body api.BusinessService true "Business service data"
func (h BusinessServiceHandler) Create(ctx *gin.Context) {
	r := &BusinessService{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.bindFailed(ctx, err)
		return
	}
	m := r.Model()
	result := h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	r.With(m)

	ctx.JSON(http.StatusCreated, r)
}

// Delete godoc
// @summary Delete a business service.
// @description Delete a business service.
// @tags delete
// @success 204
// @router /controls/business-service/{id} [delete]
// @param id path string true "Business service ID"
func (h BusinessServiceHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.BusinessService{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a business service.
// @description Update a business service.
// @tags update
// @accept json
// @success 204
// @router /controls/business-service/{id} [put]
// @param id path string true "Business service ID"
// @param business_service body api.BusinessService true "Business service data"
func (h BusinessServiceHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	r := &BusinessService{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.bindFailed(ctx, err)
		return
	}
	updates := r.Model()
	result := h.DB.Model(&model.BusinessService{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// BusinessService REST resource.
type BusinessService struct {
	Resource
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Owner       struct {
		ID          *uint  `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"owner"`
}

//
// With updates the resource with the model.
func (r *BusinessService) With(m *model.BusinessService) {
	r.Resource.With(&m.Model)
	r.Name = m.Name
	r.Description = m.Description
	r.Owner.ID = m.OwnerID
	if m.Owner != nil {
		r.Owner.DisplayName = m.Owner.DisplayName
	}
}

//
// Model builds a model.
func (r *BusinessService) Model() (m *model.BusinessService) {
	m = &model.BusinessService{
		Name:        r.Name,
		Description: r.Description,
		OwnerID:     r.Owner.ID,
	}
	m.ID = r.ID
	return
}
