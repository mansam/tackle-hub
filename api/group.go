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
	StakeholderGroupKind = "stakeholder-group"
)

//
// Routes
const (
	StakeholderGroupsRoot = ControlsRoot + "/stakeholder-group"
	StakeholderGroupRoot  = StakeholderGroupsRoot + "/:" + ID
)

//
// StakeholderGroupHandler handles stakeholder-group routes.
type StakeholderGroupHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h StakeholderGroupHandler) AddRoutes(e *gin.Engine) {
	e.GET(StakeholderGroupsRoot, h.List)
	e.GET(StakeholderGroupsRoot+"/", h.List)
	e.POST(StakeholderGroupsRoot, h.Create)
	e.GET(StakeholderGroupRoot, h.Get)
	e.PUT(StakeholderGroupRoot, h.Update)
	e.DELETE(StakeholderGroupRoot, h.Delete)
}

// Get godoc
// @summary Get a stakeholder group by ID.
// @description Get a stakeholder group by ID.
// @tags get
// @produce json
// @success 200 {object} api.StakeholderGroup
// @router /controls/stakeholder-group/:id [get]
// @param id path string true "Stakeholder Group ID"
func (h StakeholderGroupHandler) Get(ctx *gin.Context) {
	m := &model.StakeholderGroup{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Stakeholders")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := StakeholderGroup{}
	r.With(m)

	ctx.JSON(http.StatusOK, m)
}

// List godoc
// @summary List all stakeholder groups.
// @description List all stakeholder groups.
// @tags get
// @produce json
// @success 200 {object} []api.StakeholderGroup
// @router /controls/stakeholder-group [get]
func (h StakeholderGroupHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.StakeholderGroup
	h.DB.Model(model.StakeholderGroup{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(h.DB, "Stakeholders")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []StakeholderGroup{}
	for i := range models {
		r := StakeholderGroup{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, StakeholderGroupKind, resources, int(count))
}

// Create godoc
// @summary Create a stakeholder group.
// @description Create a stakeholder group.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.StakeholderGroup
// @router /controls/stakeholder-group [post]
// @param stakeholder_group body api.StakeholderGroup true "Stakeholder Group data"
func (h StakeholderGroupHandler) Create(ctx *gin.Context) {
	r := &StakeholderGroup{}
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
	r.With(m)

	ctx.JSON(http.StatusCreated, r)
}

// Delete godoc
// @summary Delete a stakeholder group.
// @description Delete a stakeholder group.
// @tags delete
// @success 204
// @router /controls/stakeholder-group/:id [delete]
// @param id path string true "Stakeholder Group ID"
func (h StakeholderGroupHandler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param(ID))
	model := &model.StakeholderGroup{}
	model.ID = uint(id)
	result := h.DB.Select("Stakeholders").Delete(model)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a stakeholder group.
// @description Update a stakeholder group.
// @tags update
// @accept json
// @success 204
// @router /controls/stakeholder-group/:id [put]
// @param id path string true "Stakeholder Group ID"
// @param stakeholder_group body api.StakeholderGroup true "Stakeholder Group data"
func (h StakeholderGroupHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	r := &StakeholderGroup{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	m := r.Model()
	result := h.DB.Model(&model.StakeholderGroup{}).Where("id = ?", id).Omit("id").Updates(m)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}
	err = h.DB.Model(&m).Association("Stakeholders").Replace("Stakeholders", m.Stakeholders)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// StakeholderGroup REST resource.
type StakeholderGroup struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name" binding:"required"`
	Description  string        `json:"description"`
	Stakeholders []Stakeholder `json:"stakeholders"`
}

//
// With updates the resource with the model.
func (r *StakeholderGroup) With(m *model.StakeholderGroup) {
	r.ID = m.ID
	r.Name = m.Name
	r.Description = m.Description
	for _, s := range m.Stakeholders {
		r.Stakeholders = append(r.Stakeholders, Stakeholder{
			ID:          s.ID,
			DisplayName: s.DisplayName,
		})
	}
}

//
// Model builds a model.
func (r *StakeholderGroup) Model() (m *model.StakeholderGroup) {
	m = &model.StakeholderGroup{
		Name:        r.Name,
		Description: r.Description,
	}
	m.ID = r.ID
	for _, s := range r.Stakeholders {
		m.Stakeholders = append(m.Stakeholders, *s.Model())
	}
	return
}
