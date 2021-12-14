package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
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
	model := StakeholderGroup{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Stakeholders")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
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
	var models []StakeholderGroup
	h.DB.Model(StakeholderGroup{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(h.DB, "Stakeholders")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	h.listResponse(ctx, StakeholderGroupKind, models, int(count))
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
	model := StakeholderGroup{}
	err := ctx.BindJSON(&model)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(&model)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, model)
}

// Delete godoc
// @summary Delete a stakeholder group.
// @description Delete a stakeholder group.
// @tags delete
// @success 204
// @router /controls/stakeholder-group/:id [delete]
// @param id path string true "Stakeholder Group ID"
func (h StakeholderGroupHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&StakeholderGroup{}, id)
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
// @produce json
// @success 204
// @router /controls/stakeholder-group/:id [put]
// @param id path string true "Stakeholder Group ID"
// @param stakeholder_group body api.StakeholderGroup true "Stakeholder Group data"
func (h StakeholderGroupHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := StakeholderGroup{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&StakeholderGroup{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// StakeholderGroup REST resource.
type StakeholderGroup = model.StakeholderGroup
