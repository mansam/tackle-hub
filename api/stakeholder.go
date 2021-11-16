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

func (h StakeholderHandler) Get(ctx *gin.Context) {
	model := models.Stakeholder{}
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

func (h StakeholderHandler) List(ctx *gin.Context) {
	var list []models.Stakeholder
	result := h.DB.Find(&list)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h StakeholderHandler) Create(ctx *gin.Context) {
	model := models.Stakeholder{}
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

func (h StakeholderHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)

	result := h.DB.Delete(&models.Stakeholder{}, "id = ?", id)
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

func (h StakeholderHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)

	updates := models.Stakeholder{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.Stakeholder{}).Where("id = ?", id).Omit("id").Updates(updates)
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
