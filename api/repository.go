package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Routes
const (
	RepositoriesRoot    = "/repositories"
	RepositoryRoot      = RepositoriesRoot + "/:" + ID
	AppRepositoriesRoot = ApplicationRoot + RepositoriesRoot
)

//
// RepositoryHandler handles repository resource routes.
type RepositoryHandler struct {
	BaseHandler
}

func (h RepositoryHandler) AddRoutes(e *gin.Engine) {
	e.GET(RepositoriesRoot, h.List)
	e.GET(RepositoriesRoot+"/", h.List)
	e.POST(RepositoriesRoot, h.Create)
	e.POST(AppRepositoriesRoot, h.CreateForApplication)
	e.GET(AppRepositoriesRoot, h.ListByApplication)
	e.GET(RepositoryRoot, h.Get)
	e.PUT(RepositoryRoot, h.Update)
	e.DELETE(RepositoryRoot, h.Delete)
}

// Get godoc
// @summary Get a repository by ID.
// @description Get a repository by ID.
// @tags get
// @produce json
// @success 200 {object} Repository
// @router /repositories/:id [get]
// @param id path string true "Repository ID"
func (h RepositoryHandler) Get(ctx *gin.Context) {
	repository := Repository{}
	id := ctx.Param(ID)
	result := h.DB.First(&repository, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, repository)
}

// List godoc
// @summary List all repositories.
// @description List all repositories.
// @tags get
// @produce json
// @success 200 {object} []Repository
// @router /repositories [get]
func (h RepositoryHandler) List(ctx *gin.Context) {
	var list []Repository
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// ListByApplication godoc
// @summary List all repositories that belong to an application.
// @description List all repositories that belong to an application.
// @tags get
// @produce json
// @success 200 {object} []Repository
// @router /application-inventory/application/:id/repositories [get]
func (h RepositoryHandler) ListByApplication(ctx *gin.Context) {
	var list []Repository
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
// @summary Create a repository.
// @description Create a repository.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Repository
// @router /repositories [post]
// @param repo body Repository true "Repository data"
func (h RepositoryHandler) Create(ctx *gin.Context) {
	repository := &Repository{}
	err := ctx.BindJSON(repository)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(repository)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, repository)
}

// CreateForApplication godoc
// @summary Create a repository for an application.
// @description Create a repository for an application.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Repository
// @router /application-inventory/application/:id/repositories [post]
// @param repo body Repository true "Repository data"
func (h RepositoryHandler) CreateForApplication(ctx *gin.Context) {
	repository := &Repository{}
	err := ctx.BindJSON(repository)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	appID := ctx.Param(ID)
	application := &model.Application{}
	result := h.DB.First(application, appID)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	repository.ApplicationID = application.ID
	result = h.DB.Create(repository)
	if result.Error != nil {
		h.createFailed(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, repository)
}

// Delete godoc
// @summary Delete a repository.
// @description Delete a repository.
// @tags delete
// @success 204
// @router /repositories/:id [delete]
// @param id path string true "Repository ID"
func (h RepositoryHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	repository := &Repository{}
	result := h.DB.First(repository, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}
	result = h.DB.Delete(repository, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a repository.
// @description Update a repository.
// @tags update
// @accept json
// @success 204
// @router /repositories/:id [put]
// @param id path string true "Repository ID"
// @param repo body Repository true "Repository data"
func (h RepositoryHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := Repository{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	db := h.DB.Model(&Repository{})
	db = db.Where("id", id)
	db = db.Omit("id", "location")
	result := db.Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

//
// Repository REST resource.
type Repository = model.Repository
