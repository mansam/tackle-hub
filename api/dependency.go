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

type DependencyHandler struct {
	BaseHandler
}

func (h DependencyHandler) AddRoutes(e *gin.Engine) {
	e.GET(DependenciesRoot, h.List)
	e.GET(DependenciesRoot+"/", h.List)
	e.POST(DependenciesRoot, h.Create)
	e.GET(DependencyRoot, h.Get)
	e.DELETE(DependencyRoot, h.Delete)
}

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

func (h DependencyHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Dependency

	to := ctx.Query("to.id")
	from := ctx.Query("from.id")
	if to == "" && from == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "'from.id' or 'to.id' must be provided"})
		return
	}

	db := h.DB
	if to != "" {
		db = db.Where("to_id = ?", to)
	} else {
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

	list := List{}
	list.With(DependencyKind, resources, int(count))
	h.hal(ctx, http.StatusOK, list)
}

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

func (h DependencyHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Dependency{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

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

func (r *Dependency) With(m *model.Dependency) {
	r.ID = m.ID
	r.To.ID = m.ToID
	r.To.Name = m.To.Name
	r.From.ID = m.FromID
	r.From.Name = m.From.Name
}

func (r *Dependency) Model() (m *model.Dependency) {
	m = &model.Dependency{
		ToID:   r.To.ID,
		FromID: r.From.ID,
	}
	return
}
