package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	ReviewKind = "review"
)

//
// Routes
const (
	ReviewsRoot = InventoryRoot + "/review"
	ReviewRoot  = ReviewsRoot + "/:" + ID
	BulkRoot    = ReviewsRoot + "/bulk"
)

//
// ReviewHandler handles review routes.
type ReviewHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h ReviewHandler) AddRoutes(e *gin.Engine) {
	e.GET(ReviewsRoot, h.List)
	e.GET(ReviewsRoot+"/", h.List)
	e.POST(ReviewsRoot, h.Create)
	e.GET(ReviewRoot, h.Get)
	e.PUT(ReviewRoot, h.Update)
	e.DELETE(ReviewRoot, h.Delete)
	e.POST(BulkRoot, h.CopyReview)
}

// Get godoc
// @summary Get a review by ID.
// @description Get a review by ID.
// @tags get
// @produce json
// @success 200 {object} []api.Review
// @router /application-inventory/review/:id [get]
// @param id path string true "Review ID"
func (h ReviewHandler) Get(ctx *gin.Context) {
	m := &model.Review{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Application")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Review{}
	r.With(m)

	ctx.JSON(http.StatusOK, r)
}

// List godoc
// @summary List all reviews.
// @description List all reviews.
// @tags get
// @produce json
// @success 200 {object} []api.Review
// @router /application-inventory/review [get]
func (h ReviewHandler) List(ctx *gin.Context) {
	var count int64
	var list []model.Review
	h.DB.Model(&model.Review{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Application")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Review{}
	for i := range list {
		r := Review{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, ReviewKind, resources, int(count))
}

// Create godoc
// @summary Create a review.
// @description Create a review.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Review
// @router /application-inventory/review [post]
// @param review body api.Review true "Review data"
func (h ReviewHandler) Create(ctx *gin.Context) {
	review := Review{}
	err := ctx.BindJSON(&review)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	m := review.Model()
	result := h.DB.Find(&model.Review{}, "applicationid", m.ApplicationID)
	if result.RowsAffected > 0 {
		h.conflict(ctx, "applicationid")
		return
	}
	result = h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	review.With(m)

	ctx.JSON(http.StatusCreated, review)
}

// Delete godoc
// @summary Delete a review.
// @description Delete a review.
// @tags delete
// @success 204
// @router /application-inventory/review/:id [delete]
// @param id path string true "Review ID"
func (h ReviewHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&Review{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a review.
// @description Update a review.
// @tags update
// @accept json
// @success 204
// @router /application-inventory/review/:id [put]
// @param id path string true "Review ID"
// @param review body api.Review true "Review data"
func (h ReviewHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := Review{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&Review{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// CopyReview godoc
// @summary Copy a review from one application to others.
// @description Copy a review from one application to others.
// @tags copy
// @accept json
// @success 204
// @router /application-inventory/review/bulk [post]
// @param copy_request body api.CopyRequest true "Review copy request data"
func (h ReviewHandler) CopyReview(ctx *gin.Context) {
	c := CopyRequest{}
	err := ctx.BindJSON(&c)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}

	m := model.Review{}
	result := h.DB.First(&m, c.SourceReview)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	for _, id := range c.TargetApplications {
		copied := model.Review{
			BusinessCriticality: m.BusinessCriticality,
			EffortEstimate:      m.EffortEstimate,
			ProposedAction:      m.ProposedAction,
			WorkPriority:        m.WorkPriority,
			Comments:            m.Comments,
			ApplicationID:       id,
		}
		existing := []model.Review{}
		result = h.DB.Find(&existing, "applicationid = ?", id)
		if result.Error != nil {
			h.createFailed(ctx, result.Error)
			return
		}
		// if the application doesn't already have a review, create one.
		if len(existing) == 0 {
			result = h.DB.Create(&copied)
			if result.Error != nil {
				h.createFailed(ctx, result.Error)
				return
			}
			// if the application already has a review, replace it with the copied review.
		} else {
			result = h.DB.Model(&model.Review{}).Where("id = ?", existing[0].ID).Updates(&copied)
			if result.Error != nil {
				h.createFailed(ctx, result.Error)
				return
			}
		}
	}
	ctx.Status(http.StatusNoContent)
}

//
// Review REST resource.
type Review struct {
	ID                  uint   `json:"id"`
	BusinessCriticality uint   `json:"businessCriticality"`
	EffortEstimate      string `json:"effortEstimate"`
	ProposedAction      string `json:"proposedAction"`
	WorkPriority        uint   `json:"workPriority"`
	Comments            string `json:"comments"`
	Application         *struct {
		ID uint `json:"id" binding:"required"`
	} `json:"application" binding:"required"`
}

// With updates the resource with the model.
func (r *Review) With(m *model.Review) {
	r.ID = m.ID
	r.BusinessCriticality = m.BusinessCriticality
	r.EffortEstimate = m.EffortEstimate
	r.ProposedAction = m.ProposedAction
	r.WorkPriority = m.WorkPriority
	r.Comments = m.Comments
	r.Application = &struct {
		ID uint `json:"id" binding:"required"`
	}{
		ID: m.ApplicationID,
	}
}

//
// Model builds a model.
func (r *Review) Model() (m *model.Review) {
	m = &model.Review{
		BusinessCriticality: r.BusinessCriticality,
		EffortEstimate:      r.EffortEstimate,
		ProposedAction:      r.ProposedAction,
		WorkPriority:        r.WorkPriority,
		Comments:            r.Comments,
	}
	if r.Application != nil {
		m.ApplicationID = r.Application.ID
	}
	m.ID = r.ID
	return
}

//
// CopyRequest REST resource.
type CopyRequest struct {
	SourceReview       uint   `json:"sourceReview" binding:"required"`
	TargetApplications []uint `json:"targetApplications" binding:"required"`
}
