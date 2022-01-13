package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Routes
const (
	ProxiesRoot = "/proxies"
	ProxyRoot   = ProxiesRoot + "/" + ID
)

//
// ProxyHandler handles proxy resource routes.
type ProxyHandler struct {
	BaseHandler
}

func (h ProxyHandler) AddRoutes(e *gin.Engine) {
	e.GET(ProxiesRoot, h.List)
	e.GET(ProxiesRoot+"/", h.List)
	e.POST(ProxiesRoot, h.Create)
	e.GET(ProxyRoot, h.Get)
	e.PUT(ProxyRoot, h.Update)
	e.DELETE(ProxyRoot, h.Delete)
}

// Get godoc
// @summary Get an proxy by ID.
// @description Get an proxy by ID.
// @tags get
// @produce json
// @success 200 {object} Proxy
// @router /proxies/:id [get]
// @param id path string true "Proxy ID"
func (h ProxyHandler) Get(ctx *gin.Context) {
	proxy := Proxy{}
	id := ctx.Param(ID)
	result := h.DB.First(&proxy, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, proxy)
}

// List godoc
// @summary List all proxies.
// @description List all proxies.
// @tags get
// @produce json
// @success 200 {object} []Proxy
// @router /proxies [get]
func (h ProxyHandler) List(ctx *gin.Context) {
	var list []Proxy
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	kind := ctx.Query("kind")
	if kind != "" {
		db = db.Where("kind", kind)
	}
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create an proxy.
// @description Create an proxy.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Proxy
// @router /proxies [post]
// @param proxy body Proxy true "Proxy data"
func (h ProxyHandler) Create(ctx *gin.Context) {
	proxy := &Proxy{}
	err := ctx.BindJSON(proxy)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(proxy)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, proxy)
}

// Delete godoc
// @summary Delete an proxy.
// @description Delete an proxy.
// @tags delete
// @success 204
// @router /proxies/:id [delete]
// @param id path string true "Proxy ID"
func (h ProxyHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	proxy := &Proxy{}
	result := h.DB.First(proxy, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}
	result = h.DB.Delete(proxy, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update an proxy.
// @description Update an proxy.
// @tags update
// @accept json
// @success 204
// @router /proxies/:id [put]
// @param id path string true "Proxy ID"
// @param proxy body Proxy true "Proxy data"
func (h ProxyHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := Proxy{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	db := h.DB.Model(&Proxy{})
	db = db.Where("id", id)
	db = db.Omit("id")
	result := db.Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Proxy REST resource.
type Proxy = model.Proxy
