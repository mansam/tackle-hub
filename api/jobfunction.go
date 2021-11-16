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
	JobFunctionsRoot = ControlsRoot + "/job_functions"
	JobFunctionRoot  = JobFunctionsRoot + "/:" + ID
)

type JobFunctionHandler struct {
	BaseHandler
}

func (h JobFunctionHandler) AddRoutes(e *gin.Engine) {
	e.GET(JobFunctionsRoot, h.List)
	e.GET(JobFunctionsRoot+"/", h.List)
	e.POST(JobFunctionsRoot, h.Create)
	e.GET(JobFunctionRoot, h.Get)
	e.PUT(JobFunctionRoot, h.Update)
	e.DELETE(JobFunctionRoot, h.Delete)
}

func (h JobFunctionHandler) Get(ctx *gin.Context) {
	model := models.JobFunction{}
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

func (h JobFunctionHandler) List(ctx *gin.Context) {
	var list []models.JobFunction
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

func (h JobFunctionHandler) Create(ctx *gin.Context) {
	model := models.JobFunction{}
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

func (h JobFunctionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)

	result := h.DB.Delete(&models.JobFunction{}, "id = ?", id)
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

func (h JobFunctionHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)

	updates := models.JobFunction{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := h.DB.Model(&models.JobFunction{}).Where("id = ?", id).Omit("id").Updates(updates)
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
