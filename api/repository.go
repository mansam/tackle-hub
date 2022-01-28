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
	e.GET(RepositoryRoot, h.Get)
	e.PUT(RepositoryRoot, h.Update)
	e.DELETE(RepositoryRoot, h.Delete)
	e.POST(AppRepositoriesRoot, h.CreateForApplication)
	e.GET(AppRepositoriesRoot, h.ListByApplication)
}

// Get godoc
// @summary Get a repository by ID.
// @description Get a repository by ID.
// @tags get
// @produce json
// @success 200 {object} Repository
// @router /repositories/{id} [get]
// @param id path string true "Repository ID"
func (h RepositoryHandler) Get(ctx *gin.Context) {
	repository := &model.Repository{}
	id := ctx.Param(ID)
	result := h.DB.First(repository, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Repository{}
	r.With(repository)

	ctx.JSON(http.StatusOK, r)
}

// List godoc
// @summary List all repositories.
// @description List all repositories.
// @tags get
// @produce json
// @success 200 {object} []Repository
// @router /repositories [get]
func (h RepositoryHandler) List(ctx *gin.Context) {
	var list []model.Repository
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Repository{}
	for i := range list {
		r := Repository{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	ctx.JSON(http.StatusOK, resources)
}

// ListByApplication godoc
// @summary List all repositories that belong to an application.
// @description List all repositories that belong to an application.
// @tags get
// @produce json
// @success 200 {object} []Repository
// @router /application-inventory/application/{id}/repositories [get]
// @param id path int true "Application ID"
func (h RepositoryHandler) ListByApplication(ctx *gin.Context) {
	var list []model.Repository
	appId := ctx.Param(ID)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = db.Where("applicationid", appId)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Repository{}
	for i := range list {
		r := Repository{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	ctx.JSON(http.StatusOK, resources)
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
		return
	}
	m := repository.Model()
	result := h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	repository.With(m)

	ctx.JSON(http.StatusCreated, repository)
}

// CreateForApplication godoc
// @summary Create a repository for an application.
// @description Create a repository for an application.
// @tags create
// @accept json
// @produce json
// @success 201 {object} Repository
// @router /application-inventory/application/{id}/repositories [post]
// @param id path int true "Application ID"
// @param repo body Repository true "Repository data"
func (h RepositoryHandler) CreateForApplication(ctx *gin.Context) {
	repository := &Repository{}
	err := ctx.BindJSON(repository)
	if err != nil {
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
	m := repository.Model()
	result = h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, err)
		return
	}
	repository.With(m)

	ctx.JSON(http.StatusCreated, repository)
}

// Delete godoc
// @summary Delete a repository.
// @description Delete a repository.
// @tags delete
// @success 204
// @router /repositories/{id} [delete]
// @param id path string true "Repository ID"
func (h RepositoryHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	repository := &model.Repository{}
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

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a repository.
// @description Update a repository.
// @tags update
// @accept json
// @success 204
// @router /repositories/{id} [put]
// @param id path string true "Repository ID"
// @param repo body Repository true "Repository data"
func (h RepositoryHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	r := &Repository{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.bindFailed(ctx, err)
		return
	}
	updates := r.Model()
	db := h.DB.Model(&model.Repository{})
	db = db.Where("id", id)
	db = db.Omit("id", "location")
	result := db.Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Repository REST resource.
type Repository struct {
	Resource
	Kind          string `json:"kind" binding:"required"`
	URL           string `json:"url" binding:"url"`
	Branch        string `json:"branch"`
	Tag           string `json:"tag"`
	Path          string `json:"path"`
	ApplicationID uint   `json:"application"`
}

//
// With updates the resource with the model.
func (r *Repository) With(m *model.Repository) {
	r.Resource.With(&m.Model)
	r.Kind = m.Kind
	r.URL = m.URL
	r.Branch = m.Branch
	r.Tag = m.Tag
	r.Path = m.Path
	r.ApplicationID = m.ApplicationID
}

//
// Model builds a model.
func (r *Repository) Model() (m *model.Repository) {
	m = &model.Repository{
		Kind:          r.Kind,
		URL:           r.URL,
		Branch:        r.Branch,
		Tag:           r.Tag,
		Path:          r.Path,
		ApplicationID: r.ApplicationID,
	}
	m.ID = r.ID

	return
}
