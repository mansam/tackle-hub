package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
	"path"
)

//
// Routes
const (
	ArtifactsRoot            = "/artifacts"
	ArtifactRoot             = ArtifactsRoot + "/:" + ID
	ApplicationArtifactsRoot = ApplicationRoot + ArtifactsRoot
)

//
// ArtifactHandler handles artifact routes.
type ArtifactHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h ArtifactHandler) AddRoutes(e *gin.Engine) {
	e.GET(ArtifactsRoot, h.List)
	e.GET(ArtifactsRoot+"/", h.List)
	e.GET(ApplicationArtifactsRoot, h.List)
	e.GET(ApplicationArtifactsRoot+"/", h.List)
	e.POST(ArtifactsRoot, h.Create)
	e.GET(ArtifactRoot, h.Get)
	e.PUT(ArtifactRoot, h.Update)
	e.DELETE(ArtifactRoot, h.Delete)
}

// Get godoc
// @summary Get a artifact by ID.
// @description Get a artifact by ID.
// @artifacts get
// @produce json
// @success 200 {object} model.Artifact
// @router /controls/artifact/:id [get]
// @param id path string true "Artifact ID"
func (h ArtifactHandler) Get(ctx *gin.Context) {
	artifact := model.Artifact{}
	id := ctx.Param(ID)
	result := h.DB.First(&artifact, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, artifact)
}

// List godoc
// @summary List all artifacts.
// @description List all artifacts.
// @artifacts get
// @produce json
// @success 200 {object} model.Artifact
// @router /controls/artifact [get]
func (h ArtifactHandler) List(ctx *gin.Context) {
	appId := ctx.Param(ID)
	var list []model.Artifact
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	if len(appId) > 0 {
		db = db.Where("application_id", appId)
	}
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a artifact.
// @description Create a artifact.
// @artifacts create
// @accept json
// @produce json
// @success 200 {object} model.Artifact
// @router /controls/artifact [post]
// @param artifact body model.Artifact true "Artifact data"
func (h ArtifactHandler) Create(ctx *gin.Context) {
	artifact := model.Artifact{}
	err := ctx.BindJSON(&artifact)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	uid := uuid.New()
	artifact.Location = path.Join(
		"/var/tackle/artifact",
		artifact.Name,
		uid.String())
	// Store uploaded file.
	// This could be filesystem (mounted PV), cloud etc.
	result := h.DB.Create(&artifact)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, artifact)
}

// Delete godoc
// @summary Delete a artifact.
// @description Delete a artifact.
// @artifacts delete
// @success 200 {object} model.Artifact
// @router /controls/artifact/:id [delete]
// @param id path string true "Artifact ID"
func (h ArtifactHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Artifact{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a artifact.
// @description Update a artifact.
// @artifacts update
// @accept json
// @produce json
// @success 200 {object} model.Artifact
// @router /controls/artifact/:id [put]
// @param id path string true "Artifact ID"
// @param artifact body model.Artifact true "Artifact data"
func (h ArtifactHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := model.Artifact{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&model.Artifact{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
