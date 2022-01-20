package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Routes
const (
	IdentitiesRoot    = "/identities"
	IdentityRoot      = IdentitiesRoot + "/:" + ID
	AppIdentitiesRoot = ApplicationRoot + IdentitiesRoot
)

//
// IdentityHandler handles identity resource routes.
type IdentityHandler struct {
	BaseHandler
}

func (h IdentityHandler) AddRoutes(e *gin.Engine) {
	e.GET(IdentitiesRoot, h.List)
	e.GET(IdentitiesRoot+"/", h.List)
	e.POST(IdentitiesRoot, h.Create)
	e.GET(IdentityRoot, h.Get)
	e.PUT(IdentityRoot, h.Update)
	e.DELETE(IdentityRoot, h.Delete)
	e.POST(AppIdentitiesRoot, h.CreateForApplication)
	e.GET(AppIdentitiesRoot, h.ListByApplication)
	e.GET(AppIdentitiesRoot+"/", h.ListByApplication)
}

// Get godoc
// @summary Get an identity by ID.
// @description Get an identity by ID.
// @tags get
// @produce json
// @success 200 {object} Identity
// @router /identities/:id [get]
// @param id path string true "Identity ID"
func (h IdentityHandler) Get(ctx *gin.Context) {
	identity := Identity{}
	id := ctx.Param(ID)
	result := h.DB.First(&identity, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, identity)
}

// List godoc
// @summary List all identities.
// @description List all identities.
// @tags get
// @produce json
// @success 200 {object} []Identity
// @router /identities [get]
func (h IdentityHandler) List(ctx *gin.Context) {
	var list []Identity
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// ListByApplication  godoc
// @summary List identities for an application.
// @description List identities for an application.
// @tags get
// @produce json
// @success 200 {object} []Identity
// @router /identities [get]
func (h IdentityHandler) ListByApplication(ctx *gin.Context) {
	var list []Identity
	appId := ctx.Param(ID)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = db.Where("application_id", appId)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create an identity.
// @description Create an identity.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Identity
// @router /identities [post]
// @param identity body Identity true "Identity data"
func (h IdentityHandler) Create(ctx *gin.Context) {
	identity := &Identity{}
	err := ctx.BindJSON(identity)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(identity)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, identity)
}

// CreateForApplication godoc
// @summary Create an identity for an application.
// @description Create an identity for an application.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Identity
// @router /identities [post]
// @param identity body Identity true "Identity data"
func (h IdentityHandler) CreateForApplication(ctx *gin.Context) {
	identity := &Identity{}
	err := ctx.BindJSON(identity)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	appId := ctx.Param(ID)
	application := &model.Application{}
	result := h.DB.First(application, appId)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	identity.ApplicationID = application.ID
	result = h.DB.Create(identity)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, identity)
}

// Delete godoc
// @summary Delete an identity.
// @description Delete an identity.
// @tags delete
// @success 204
// @router /identities/:id [delete]
// @param id path string true "Identity ID"
func (h IdentityHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	identity := &Identity{}
	result := h.DB.First(identity, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}
	result = h.DB.Delete(identity, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update an identity.
// @description Update an identity.
// @tags update
// @accept json
// @success 204
// @router /identities/:id [put]
// @param id path string true "Identity ID"
// @param identity body Identity true "Identity data"
func (h IdentityHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := Identity{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	db := h.DB.Model(&Identity{})
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
// Identity REST resource.
type Identity = model.Identity
