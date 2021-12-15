package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	TagTypeKind = "tag-type"
)

//
// Routes
const (
	TagTypesRoot = ControlsRoot + "/tag-type"
	TagTypeRoot  = TagTypesRoot + "/:" + ID
)

//
// TagTypeHandler handles the tag-type route.
type TagTypeHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h TagTypeHandler) AddRoutes(e *gin.Engine) {
	e.GET(TagTypesRoot, h.List)
	e.GET(TagTypesRoot+"/", h.List)
	e.POST(TagTypesRoot, h.Create)
	e.GET(TagTypeRoot, h.Get)
	e.PUT(TagTypeRoot, h.Update)
	e.DELETE(TagTypeRoot, h.Delete)
}

// Get godoc
// @summary Get a tag type by ID.
// @description Get a tag type by ID.
// @tags get
// @produce json
// @success 200 {object} api.TagType
// @router /controls/tag-type/:id [get]
// @param id path string true "Tag Type ID"
func (h TagTypeHandler) Get(ctx *gin.Context) {
	model := TagType{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Tags")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all tag types.
// @description List all tag types.
// @tags get
// @produce json
// @success 200 {object} []api.TagType
// @router /controls/tag-type [get]
func (h TagTypeHandler) List(ctx *gin.Context) {
	var count int64
	var models []TagType
	h.DB.Model(TagType{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Tags")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	h.listResponse(ctx, TagTypeKind, models, int(count))
}

// Create godoc
// @summary Create a tag type.
// @description Create a tag type.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.TagType
// @router /controls/tag-type [post]
// @param tag_type body api.TagType true "Tag Type data"
func (h TagTypeHandler) Create(ctx *gin.Context) {
	model := TagType{}
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
// @summary Delete a tag type.
// @description Delete a tag type.
// @tags delete
// @success 204
// @router /controls/tag-type/:id [delete]
// @param id path string true "Tag Type ID"
func (h TagTypeHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&TagType{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a tag type.
// @description Update a tag type.
// @tags update
// @accept json
// @success 204
// @router /controls/tag-type/:id [put]
// @param id path string true "Tag Type ID"
// @param tag_type body api.TagType true "Tag Type data"
func (h TagTypeHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := TagType{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&TagType{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// TagType REST resource.
type TagType = model.TagType
