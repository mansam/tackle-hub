package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Routes
const (
	GroupsRoot = ControlsRoot + "/stakeholder-group"
	GroupRoot  = GroupsRoot + "/:" + ID
)

//
// StakeholderGroupHandler handles stakeholder-group routes.
type StakeholderGroupHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h StakeholderGroupHandler) AddRoutes(e *gin.Engine) {
	e.GET(GroupsRoot, h.List)
	e.GET(GroupsRoot+"/", h.List)
	e.POST(GroupsRoot, h.Create)
	e.GET(GroupRoot, h.Get)
	e.PUT(GroupRoot, h.Update)
	e.DELETE(GroupRoot, h.Delete)
}

// Get godoc
// @summary Get a stakeholder group by ID.
// @description Get a stakeholder group by ID.
// @tags get
// @produce json
// @success 200 {object} models.StakeholderGroup
// @router /controls/stakeholder-group/:id [get]
// @param id path string true "Stakeholder Group ID"
func (h StakeholderGroupHandler) Get(ctx *gin.Context) {
	model := model.StakeholderGroup{}
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
// @success 200 {object} models.StakeholderGroup
// @router /controls/stakeholder-group [get]
func (h StakeholderGroupHandler) List(ctx *gin.Context) {
	var list []model.StakeholderGroup
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(h.DB, "Stakeholders")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a stakeholder group.
// @description Create a stakeholder group.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.StakeholderGroup
// @router /controls/stakeholder-group [post]
// @param stakeholder_group body models.StakeholderGroup true "Stakeholder Group data"
func (h StakeholderGroupHandler) Create(ctx *gin.Context) {
	model := model.StakeholderGroup{}
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

	ctx.JSON(http.StatusOK, model)
}

// Delete godoc
// @summary Delete a stakeholder group.
// @description Delete a stakeholder group.
// @tags delete
// @success 200 {object} models.StakeholderGroup
// @router /controls/stakeholder-group/:id [delete]
// @param id path string true "Stakeholder Group ID"
func (h StakeholderGroupHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.StakeholderGroup{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a stakeholder group.
// @description Update a stakeholder group.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.StakeholderGroup
// @router /controls/stakeholder-group/:id [put]
// @param id path string true "Stakeholder Group ID"
// @param stakeholder_group body models.StakeholderGroup true "Stakeholder Group data"
func (h StakeholderGroupHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := model.StakeholderGroup{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&model.StakeholderGroup{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
