package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	StakeholderKind = "stakeholder"
)

//
// Routes
const (
	StakeholdersRoot = ControlsRoot + "/stakeholder"
	StakeholderRoot  = StakeholdersRoot + "/:" + ID
)

//
// StakeholderHandler handles stakeholder routes.
type StakeholderHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h StakeholderHandler) AddRoutes(e *gin.Engine) {
	e.GET(StakeholdersRoot, h.List)
	e.GET(StakeholdersRoot+"/", h.List)
	e.POST(StakeholdersRoot, h.Create)
	e.GET(StakeholderRoot, h.Get)
	e.PUT(StakeholderRoot, h.Update)
	e.DELETE(StakeholderRoot, h.Delete)
}

// Get godoc
// @summary Get a stakeholder by ID.
// @description Get a stakeholder by ID.
// @tags get
// @produce json
// @success 200 {object} api.Stakeholder
// @router /controls/stakeholder/:id [get]
// @param id path string true "Stakeholder ID"
func (h StakeholderHandler) Get(ctx *gin.Context) {
	model := model.Stakeholder{}
	id := ctx.Param(ID)
	db := h.preLoad(
		h.DB,
		"JobFunction",
		"BusinessServices",
		"StakeholderGroups")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	resource := Stakeholder{}
	resource.With(&model)
	ctx.JSON(http.StatusOK, resource)
}

// List godoc
// @summary List all stakeholders.
// @description List all stakeholders.
// @tags get
// @produce json
// @success 200 {object} []api.Stakeholder
// @router /controls/stakeholder [get]
func (h StakeholderHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Stakeholder
	h.DB.Model(model.Stakeholder{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(
		db,
		"JobFunction",
		"BusinessServices",
		"Groups")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Stakeholder{}
	for i := range models {
		r := Stakeholder{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, StakeholderKind, resources, int(count))
}

// Create godoc
// @summary Create a stakeholder.
// @description Create a stakeholder.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Stakeholder
// @router /controls/stakeholder [post]
// @param stakeholder body api.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Create(ctx *gin.Context) {
	resource := Stakeholder{}
	err := ctx.BindJSON(&resource)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	model := resource.Model()
	result := h.DB.Create(model)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	resource.With(model)
	ctx.JSON(http.StatusCreated, resource)
}

// Delete godoc
// @summary Delete a stakeholder.
// @description Delete a stakeholder.
// @tags delete
// @success 204
// @router /controls/stakeholder/:id [delete]
// @param id path string true "Stakeholder ID"
func (h StakeholderHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Stakeholder{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a stakeholder.
// @description Update a stakeholder.
// @tags update
// @accept json
// @success 204
// @router /controls/stakeholder/:id [put]
// @param id path string true "Stakeholder ID"
// @param stakeholder body api.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	resource := Stakeholder{}
	err := ctx.BindJSON(&resource)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	updates := resource.Model()
	result := h.DB.Model(&model.Stakeholder{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Stakeholder REST resource.
type Stakeholder struct {
	ID               uint                     `json:"id"`
	DisplayName      string                   `json:"displayName"`
	Email            string                   `json:"email"`
	Groups           []model.StakeholderGroup `json:"stakeholderGroups"`
	BusinessServices []model.BusinessService  `json:"businessServices"`
	JobFunction      struct {
		ID   *uint  `json:"id"`
		Role string `json:"role"`
	} `json:"jobFunction"`
}

//
// With updates the resource with the model.
func (r *Stakeholder) With(m *model.Stakeholder) {
	r.ID = m.ID
	r.DisplayName = m.DisplayName
	r.Email = m.Email
	r.Groups = m.Groups
	r.BusinessServices = m.BusinessServices
	r.JobFunction.ID = m.JobFunctionID
	if m.JobFunction != nil {
		r.JobFunction.Role = m.JobFunction.Role
	}
}

//
// Model builds a model.
func (r *Stakeholder) Model() (m *model.Stakeholder) {
	m = &model.Stakeholder{
		DisplayName:      r.DisplayName,
		Email:            r.Email,
		Groups:           r.Groups,
		BusinessServices: r.BusinessServices,
		JobFunctionID:    r.JobFunction.ID,
	}
	m.ID = r.ID
	return
}
