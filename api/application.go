package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
	"strconv"
)

//
// Kind
const (
	ApplicationKind = "application"
)

//
// Routes
const (
	ApplicationsRoot = InventoryRoot + "/application"
	ApplicationRoot  = ApplicationsRoot + "/:" + ID
)

//
// ApplicationHandler handles application resource routes.
type ApplicationHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
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
// @success 200 {object} api.Application
// @router /application-inventory/application/:id [get]
// @param id path int true "Application ID"
func (h ApplicationHandler) Get(ctx *gin.Context) {
	m := &model.Application{}
	id := ctx.Param(ID)
	db := h.BaseHandler.preLoad(
		h.DB,
		"Tags",
		"Review",
		"BusinessService")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Application{}
	r.With(m)

	ctx.JSON(http.StatusOK, r)
}

// List godoc
// @summary List all applications.
// @description List all applications.
// @tags list
// @produce json
// @success 200 {object} []api.Application
// @router /application-inventory/application [get]
func (h ApplicationHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Application
	h.DB.Model(model.Application{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(
		db,
		"Tags",
		"Review",
		"BusinessService")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Application{}
	for i := range models {
		r := Application{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	list := List{}
	list.With(ApplicationKind, resources, int(count))
	h.hal(ctx, http.StatusOK, list)
}

// Create godoc
// @summary Create an application.
// @description Create an application.
// @tags create
// @accept json
// @produce json
// @success 200 {object} api.Application
// @router /application-inventory/application [post]
// @param application body model.Application true "Application data"
func (h ApplicationHandler) Create(ctx *gin.Context) {
	r := &Application{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	m := r.Model()
	result := h.DB.Create(&m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	err = h.DB.Model(m).Association("Tags").Replace("Tags", m.Tags)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	r.With(m)

	ctx.JSON(http.StatusCreated, r)
}

// Delete godoc
// @summary Delete an application.
// @description Delete an application.
// @tags delete
// @success 200
// @router /application-inventory/application/:id [delete]
// @param id path int true "Application id"
func (h ApplicationHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Application{}, id)
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
// @success 200 {object} api.Application
// @router /application-inventory/application/:id [put]
// @param id path int true "Application id"
// @param application body api.Application true "Application data"
func (h ApplicationHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	r := &Application{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	m := r.Model()
	result := h.DB.Model(&model.Application{}).Where("id = ?", id).Omit("id").Updates(m)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}
	err = h.DB.Model(m).Association("Tags").Replace("Tags", m.Tags)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

//
// Application REST resource.
type Application struct {
	model.Application
	Tags            []string `json:"tags"`
	BusinessService string   `json:"businessService"`
}

//
// With updates the resource using the model.
func (r *Application) With(m *model.Application) {
	r.Application = *m
	r.BusinessService = strconv.Itoa(int(m.BusinessServiceID))
	for _, tag := range m.Tags {
		r.Tags = append(
			r.Tags,
			strconv.Itoa(int(tag.ID)))
	}
}

//
// Model builds a model.
func (r *Application) Model() (m *model.Application) {
	m = &r.Application
	if len(r.BusinessService) > 0 {
		id, _ := strconv.Atoi(r.BusinessService)
		m.BusinessServiceID = uint(id)
	}
	for _, tagID := range r.Tags {
		id, _ := strconv.Atoi(tagID)
		m.Tags = append(
			m.Tags,
			model.Tag{
				Model: model.Model{
					ID: uint(id),
				},
			})
	}

	return
}
