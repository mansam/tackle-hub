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
	e.POST(BulkRoot, h.Copy)
}

// Get godoc
// @summary Get a review by ID.
// @description Get a review by ID.
// @tags get
// @produce json
// @success 200 {object} model.Review
// @router /application-inventory/review/:id [get]
// @param id path string true "Review ID"
func (h ReviewHandler) Get(ctx *gin.Context) {
	model := model.Review{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Application")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all reviews.
// @description List all reviews.
// @tags get
// @produce json
// @success 200 {object} model.Review
// @router /application-inventory/review [get]
func (h ReviewHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.Review
	h.DB.Model(model.Review{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Application")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	list := List{}
	list.With(ReviewKind, models, int(count))
	h.hal(ctx, http.StatusOK, list)
}

// Create godoc
// @summary Create a review.
// @description Create a review.
// @tags create
// @accept json
// @produce json
// @success 200 {object} model.Review
// @router /application-inventory/review [post]
// @param review body model.Review true "Review data"
func (h ReviewHandler) Create(ctx *gin.Context) {
	review := Review{}
	err := ctx.BindJSON(&review)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	model := review.Model()
	result := h.DB.Create(model)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, model)
}

func (h ReviewHandler) Copy(ctx *gin.Context) {
	b := Bulk{}
	err := ctx.BindJSON(&b)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	m := model.Review{}
	result := h.DB.First(&m, b.SourceReview)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	for _, id := range b.TargetApplications {
		copied := model.Review{
			BusinessCriticality: m.BusinessCriticality,
			EffortEstimate:      m.EffortEstimate,
			ProposedAction:      m.ProposedAction,
			WorkPriority:        m.WorkPriority,
			Comments:            m.Comments,
			ApplicationID:       id,
		}
		existing := []model.Review{}
		result = h.DB.Find(&existing, "application_id = ?", id)
		if result.Error != nil {
			h.createFailed(ctx, result.Error)
			return
		}
		if len(existing) == 0 {
			result = h.DB.Create(&copied)
			if result.Error != nil {
				h.createFailed(ctx, result.Error)
				return
			}
		} else {
			result = h.DB.Model(&model.Review{}).Where("id = ?", existing[0].ID).Updates(&copied)
			if result.Error != nil {
				h.createFailed(ctx, result.Error)
				return
			}
		}
	}
}

// Delete godoc
// @summary Delete a review.
// @description Delete a review.
// @tags delete
// @success 200 {object} model.Review
// @router /application-inventory/review/:id [delete]
// @param id path string true "Review ID"
func (h ReviewHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Review{}, id)
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
// @produce json
// @success 200 {object} model.Review
// @router /application-inventory/review/:id [put]
// @param id path string true "Review ID"
// @param review body model.Review true "Review data"
func (h ReviewHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	review := Review{}
	err := ctx.BindJSON(&review)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	updates := review.Model()
	result := h.DB.Model(&model.Review{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// Review request struct.
type Review struct {
	BusinessCriticality uint   `json:"businessCriticality"`
	EffortEstimate      string `json:"effortEstimate"`
	ProposedAction      string `json:"proposedAction"`
	WorkPriority        uint   `json:"workPriority"`
	Comments            string `json:"comments"`
	Application         *struct {
		ID uint `json:"id"`
	} `json:"application"`
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
	return
}

//
//
type Bulk struct {
	SourceReview       uint   `json:"sourceReview"`
	TargetApplications []uint `json:"targetApplications"`
}
