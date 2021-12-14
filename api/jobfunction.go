package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	JobFunctionKind = "job-function"
)

//
// Routes
const (
	JobFunctionsRoot = ControlsRoot + "/job-function"
	JobFunctionRoot  = JobFunctionsRoot + "/:" + ID
)

//
// JobFunctionHandler handles job-function routes.
type JobFunctionHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h JobFunctionHandler) AddRoutes(e *gin.Engine) {
	e.GET(JobFunctionsRoot, h.List)
	e.GET(JobFunctionsRoot+"/", h.List)
	e.POST(JobFunctionsRoot, h.Create)
	e.GET(JobFunctionRoot, h.Get)
	e.PUT(JobFunctionRoot, h.Update)
	e.DELETE(JobFunctionRoot, h.Delete)
}

// Get godoc
// @summary Get a job function by ID.
// @description Get a job function by ID.
// @tags get
// @produce json
// @success 200 {object} []api.JobFunction
// @router /controls/job-function/:id [get]
// @param id path string true "Job Function ID"
func (h JobFunctionHandler) Get(ctx *gin.Context) {
	model := JobFunction{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Stakeholders")
	result := db.First(&model, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// List godoc
// @summary List all job functions.
// @description List all job functions.
// @tags get
// @produce json
// @success 200 {object} []api.JobFunction
// @router /controls/job-function [get]
func (h JobFunctionHandler) List(ctx *gin.Context) {
	var count int64
	var models []JobFunction
	h.DB.Model(JobFunction{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Stakeholders")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	h.listResponse(ctx, JobFunctionKind, models, int(count))
}

// Create godoc
// @summary Create a job function.
// @description Create a job function.
// @tags create
// @accept json
// @produce json
// @success 200 {object} api.JobFunction
// @router /controls/job-function [post]
// @param job_function body api.JobFunction true "Job Function data"
func (h JobFunctionHandler) Create(ctx *gin.Context) {
	model := JobFunction{}
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
// @summary Delete a job function.
// @description Delete a job function.
// @tags delete
// @success 204
// @router /controls/job-function/:id [delete]
// @param id path string true "Job Function ID"
func (h JobFunctionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&JobFunction{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a job function.
// @description Update a job function.
// @tags update
// @accept json
// @success 204
// @router /controls/job-function/:id [put]
// @param id path string true "Job Function ID"
// @param job_function body api.JobFunction true "Job Function data"
func (h JobFunctionHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := JobFunction{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	result := h.DB.Model(&JobFunction{}).Where("id = ?", id).Omit("id").Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// JobFunction REST resrouce.
type JobFunction = model.JobFunction
