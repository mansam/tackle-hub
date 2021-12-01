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
// @success 200 {object} model.Stakeholder
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

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all stakeholders.
// @description List all stakeholders.
// @tags get
// @produce json
// @success 200 {object} model.Stakeholder
// @router /controls/stakeholder [get]
func (h StakeholderHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Stakeholder
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(
		db,
		"JobFunction",
		"BusinessServices",
		"Groups")
	result := db.Find(&models).Count(&count)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	list := List{}
	list.With(StakeholderKind, models, int(count))
	h.hal(ctx, http.StatusOK, list)
}

// Create godoc
// @summary Create a stakeholder.
// @description Create a stakeholder.
// @tags create
// @accept json
// @produce json
// @success 200 {object} model.Stakeholder
// @router /controls/stakeholder [post]
// @param stakeholder body model.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Create(ctx *gin.Context) {
	model := model.Stakeholder{}
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
// @summary Delete a stakeholder.
// @description Delete a stakeholder.
// @tags delete
// @success 200 {object} model.Stakeholder
// @router /controls/stakeholder/:id [delete]
// @param id path string true "Stakeholder ID"
func (h StakeholderHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Stakeholder{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a stakeholder.
// @description Update a stakeholder.
// @tags update
// @accept json
// @produce json
// @success 200 {object} model.Stakeholder
// @router /controls/stakeholder/:id [put]
// @param id path string true "Stakeholder ID"
// @param stakeholder body model.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := model.Stakeholder{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&model.Stakeholder{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
