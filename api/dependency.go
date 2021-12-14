package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	DependencyKind = "applications-dependency"
)

//
// Routes
const (
	DependenciesRoot = InventoryRoot + "/applications-dependency"
	DependencyRoot   = DependenciesRoot + "/:" + ID
)

//
// DependencyHandler handles application dependency routes.
type DependencyHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h DependencyHandler) AddRoutes(e *gin.Engine) {
	e.GET(DependenciesRoot, h.List)
	e.GET(DependenciesRoot+"/", h.List)
	e.POST(DependenciesRoot, h.Create)
	e.GET(DependencyRoot, h.Get)
	e.DELETE(DependencyRoot, h.Delete)
}

// Get godoc
// @summary Get a dependency by ID.
// @description Get a dependency by ID.
// @tags get
// @produce json
// @success 200 {object} api.Dependency
// @router /application-inventory/applications-dependency/:id [get]
// @param id path string true "Dependency ID"
func (h DependencyHandler) Get(ctx *gin.Context) {
	m := &model.Dependency{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "To", "From")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Dependency{}
	r.With(m)

	ctx.JSON(http.StatusOK, r)
}

//
// List godoc
// @summary List all dependencies.
// @description List all dependencies.
// @tags list
// @produce json
// @success 200 {object} []api.Dependency
// @router /application-inventory/applications-dependency [get]
func (h DependencyHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Dependency

	db := h.DB
	to := ctx.Query("to.id")
	from := ctx.Query("from.id")
	if to != "" {
		db = db.Where("to_id = ?", to)
	} else if from != "" {
		db = db.Where("from_id = ?", from)
	}

	db.Model(model.Dependency{}).Count(&count)
	pagination := NewPagination(ctx)
	db = pagination.apply(db)
	db = h.preLoad(db, "To", "From")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	resources := []Dependency{}
	for i := range models {
		r := Dependency{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, DependencyKind, resources, int(count))
}

// Create godoc
// @summary Create a dependency.
// @description Create a dependency.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Dependency
// @router /application-inventory/applications-dependency [post]
// @param applications_dependency body Dependency true "Dependency data"
func (h DependencyHandler) Create(ctx *gin.Context) {
	r := Dependency{}
	err := ctx.BindJSON(&r)
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

	ctx.JSON(http.StatusCreated, r)
}

// Delete godoc
// @summary Delete a dependency.
// @description Delete a dependency.
// @tags delete
// @accept json
// @produce json
// @success 204
// @router /application-inventory/applications-dependency/:id [delete]
// @param id path string true "Dependency id"
func (h DependencyHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Dependency{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Dependency REST resource.
type Dependency struct {
	ID uint `json:"id"`
	To struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"to"`
	From struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"from"`
}

//
// With updates the resource using the model.
func (r *Dependency) With(m *model.Dependency) {
	r.ID = m.ID
	r.To.ID = m.ToID
	r.To.Name = m.To.Name
	r.From.ID = m.FromID
	r.From.Name = m.From.Name
}

// Model builds a model.Dependency.
func (r *Dependency) Model() (m *model.Dependency) {
	m = &model.Dependency{
		ToID:   r.To.ID,
		FromID: r.From.ID,
	}
	return
}
