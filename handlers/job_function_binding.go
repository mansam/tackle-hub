package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/db"
	"github.com/konveyor/tackle-hub/models"
	"gorm.io/gorm"
	"net/http"
)

//
// Routes
const (
	JobFunctionBindingsRoot = "/job_function_binding"
	JobFunctionBindingParam = "job_function_binding"
	JobFunctionBindingRoot  = JobFunctionBindingsRoot + "/:" + JobFunctionBindingParam
)

type JobFunctionBindingHandler struct{}

func (h *JobFunctionBindingHandler) AddRoutes(e *gin.Engine) {
	e.GET(JobFunctionBindingsRoot, h.List)
	e.GET(JobFunctionBindingsRoot+"/", h.List)
	e.POST(JobFunctionBindingsRoot, h.Create)
	e.GET(JobFunctionBindingRoot, h.Get)
	e.PUT(JobFunctionBindingRoot, h.Update)
	e.DELETE(JobFunctionBindingRoot, h.Delete)
}

func (h *JobFunctionBindingHandler) Get(ctx *gin.Context) {
	model := models.JobFunctionBinding{}
	id := ctx.Param(JobFunctionBindingParam)
	result := db.DB.First(&model, "id = ?", id)
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

func (h *JobFunctionBindingHandler) List(ctx *gin.Context) {
	var list []models.JobFunctionBinding
	result := db.DB.Find(&list)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *JobFunctionBindingHandler) Create(ctx *gin.Context) {
	model := models.JobFunctionBinding{}
	err := ctx.BindJSON(&model)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := db.DB.Create(&model)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, model)
}

func (h *JobFunctionBindingHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(JobFunctionBindingParam)

	result := db.DB.Delete(&models.JobFunctionBinding{}, "id = ?", id)
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

func (h *JobFunctionBindingHandler) Update(ctx *gin.Context) {
	id := ctx.Param(JobFunctionBindingParam)

	updates := models.JobFunctionBinding{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := db.DB.Model(&models.JobFunctionBinding{}).Where("id = ?", id).Omit("id").Updates(updates)
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
