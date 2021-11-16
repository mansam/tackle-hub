package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/models"
	"gorm.io/gorm"
	"net/http"
)

//
// Routes
const (
	ReviewsRoot = InventoryRoot + "/review"
	ReviewRoot  = ReviewsRoot + "/:" + ID
)

type ReviewHandler struct {
	BaseHandler
}

func (h ReviewHandler) AddRoutes(e *gin.Engine) {
	e.GET(ReviewsRoot, h.List)
	e.GET(ReviewsRoot+"/", h.List)
	e.POST(ReviewsRoot, h.Create)
	e.GET(ReviewRoot, h.Get)
	e.PUT(ReviewRoot, h.Update)
	e.DELETE(ReviewRoot, h.Delete)
}

// Get godoc
// @summary Get a review by ID.
// @description Get a review by ID.
// @tags get
// @produce json
// @success 200 {object} models.Review
// @router /application-inventory/review/:id [get]
// @param id path string true "Review ID"
func (h ReviewHandler) Get(ctx *gin.Context) {
	model := models.Review{}
	id := ctx.Param(ID)
	result := h.DB.First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": MsgNotFound,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all reviews.
// @description List all reviews.
// @tags get
// @produce json
// @success 200 {object} models.Review
// @router /application-inventory/review [get]
func (h ReviewHandler) List(ctx *gin.Context) {
	var list []models.Review
	page := NewPagination(ctx)
	result := h.DB.Offset(page.Offset).Limit(page.Limit).Order(page.Sort).Find(&list)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// Create godoc
// @summary Create a review.
// @description Create a review.
// @tags create
// @accept json
// @produce json
// @success 200 {object} models.Review
// @router /application-inventory/review [post]
// @param review body models.Review true "Review data"
func (h ReviewHandler) Create(ctx *gin.Context) {
	model := models.Review{}
	err := ctx.BindJSON(&model)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Create(&model)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, model)
}

// Delete godoc
// @summary Delete a review.
// @description Delete a review.
// @tags delete
// @success 200 {object} models.Review
// @router /application-inventory/review/:id [delete]
// @param id path string true "Review ID"
func (h ReviewHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&models.Review{}, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusOK)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a review.
// @description Update a review.
// @tags update
// @accept json
// @produce json
// @success 200 {object} models.Review
// @router /application-inventory/review/:id [put]
// @param id path string true "Review ID"
// @param review body models.Review true "Review data"
func (h ReviewHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := models.Review{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.Review{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": MsgNotFound,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": MsgInternalServerError,
			})
			log.Error(result.Error, MsgInternalServerError)
			return
		}
	}
	ctx.Status(http.StatusOK)
}
