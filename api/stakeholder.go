package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/models"
	"net/http"
)

//
// Routes
const (
	StakeholdersRoot = ControlsRoot + "/stakeholder"
	StakeholderRoot  = StakeholdersRoot + "/:" + ID
)

type StakeholderHandler struct {
	BaseHandler
}

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
// @success 200 {object} models.Stakeholder
// @router /controls/stakeholder/:id [get]
// @param id path string true "Stakeholder ID"
func (h StakeholderHandler) Get(ctx *gin.Context) {
	model := models.Stakeholder{}
	id := ctx.Param(ID)
	result := h.DB.First(&model, "id = ?", id)
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
// @success 200 {object} models.Stakeholder
// @router /controls/stakeholder [get]
func (h StakeholderHandler) List(ctx *gin.Context) {
	var list []models.Stakeholder
	page := NewPagination(ctx)
	result := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort).Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a stakeholder.
// @description Create a stakeholder.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.Stakeholder
// @router /controls/stakeholder [post]
// @param stakeholder body models.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Create(ctx *gin.Context) {
	model := models.Stakeholder{}
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
// @summary Delete a stakeholder.
// @description Delete a stakeholder.
// @tags delete
// @success 200 {object} models.Stakeholder
// @router /controls/stakeholder/:id [delete]
// @param id path string true "Stakeholder ID"
func (h StakeholderHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&models.Stakeholder{}, "id = ?", id)
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
// @success 200 {object} models.Stakeholder
// @router /controls/stakeholder/:id [put]
// @param id path string true "Stakeholder ID"
// @param stakeholder body models.Stakeholder true "Stakeholder data"
func (h StakeholderHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := models.Stakeholder{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&models.Stakeholder{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
