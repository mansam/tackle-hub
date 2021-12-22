package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	TagKind = "tag"
)

//
// Routes
const (
	TagsRoot = ControlsRoot + "/tag"
	TagRoot  = TagsRoot + "/:" + ID
)

//
// TagHandler handles tag routes.
type TagHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h TagHandler) AddRoutes(e *gin.Engine) {
	e.GET(TagsRoot, h.List)
	e.GET(TagsRoot+"/", h.List)
	e.POST(TagsRoot, h.Create)
	e.GET(TagRoot, h.Get)
	e.PUT(TagRoot, h.Update)
	e.DELETE(TagRoot, h.Delete)
}

// Get godoc
// @summary Get a tag by ID.
// @description Get a tag by ID.
// @tags get
// @produce json
// @success 200 {object} api.Tag
// @router /controls/tag/:id [get]
// @param id path string true "Tag ID"
func (h TagHandler) Get(ctx *gin.Context) {
	model := model.Tag{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "TagType")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	resource := Tag{}
	resource.With(&model)
	ctx.JSON(http.StatusOK, resource)
}

// List godoc
// @summary List all tags.
// @description List all tags.
// @tags get
// @produce json
// @success 200 {object} []api.Tag
// @router /controls/tag [get]
func (h TagHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Tag
	h.DB.Model(model.Tag{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "TagType")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Tag{}
	for i := range models {
		r := Tag{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, TagKind, resources, int(count))
}

// Create godoc
// @summary Create a tag.
// @description Create a tag.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Tag
// @router /controls/tag [post]
// @param tag body Tag true "Tag data"
func (h TagHandler) Create(ctx *gin.Context) {
	resource := Tag{}
	err := ctx.BindJSON(&resource)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	model := resource.Model()
	result := h.DB.Create(&model)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusCreated, model)
}

// Delete godoc
// @summary Delete a tag.
// @description Delete a tag.
// @tags delete
// @success 204
// @router /controls/tag/:id [delete]
// @param id path string true "Tag ID"
func (h TagHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Tag{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a tag.
// @description Update a tag.
// @tags update
// @accept json
// @success 204
// @router /controls/tag/:id [put]
// @param id path string true "Tag ID"
// @param tag body api.Tag true "Tag data"
func (h TagHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	resource := Tag{}
	err := ctx.BindJSON(&resource)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	updates := resource.Model()
	result := h.DB.Model(&model.Tag{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Tag REST resource.
type Tag struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	TagType struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"tagType"`
}

//
// With updates the resource with the model.
func (r *Tag) With(m *model.Tag) {
	r.ID = m.ID
	r.Name = m.Name
	r.TagType.ID = m.TagTypeID
	r.TagType.Name = m.TagType.Name
}

//
// Model builds a model.
func (r *Tag) Model() (m *model.Tag) {
	m = &model.Tag{
		Name:      r.Name,
		TagTypeID: r.TagType.ID,
	}
	m.ID = r.ID
	return
}
