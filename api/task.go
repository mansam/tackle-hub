package api

import (
	"context"
	"github.com/gin-gonic/gin"
	crd "github.com/konveyor/tackle-hub/k8s/api/tackle/v1alpha1"
	"github.com/konveyor/tackle-hub/model"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"
)

//
// Routes
const (
	TasksRoot      = "/tasks"
	TaskRoot       = TasksRoot + "/:" + ID
	TaskReportRoot = TaskRoot + "/report"
	AddonTasksRoot = AddonRoot + "/tasks"
)

//
// TaskHandler handles task routes.
type TaskHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h TaskHandler) AddRoutes(e *gin.Engine) {
	e.GET(TasksRoot, h.List)
	e.GET(TasksRoot+"/", h.List)
	e.POST(TasksRoot, h.Create)
	e.GET(TaskRoot, h.Get)
	e.PUT(TaskRoot, h.Update)
	e.POST(TaskReportRoot, h.CreateReport)
	e.PUT(TaskReportRoot, h.UpdateReport)
	e.POST(AddonTasksRoot, h.AddonCreate)
	e.DELETE(TaskRoot, h.Delete)
}

// Get godoc
// @summary Get a task by ID.
// @description Get a task by ID.
// @tags get
// @produce json
// @success 200 {object} api.Task
// @router /tasks/:id [get]
// @param id path string true "Task ID"
func (h TaskHandler) Get(ctx *gin.Context) {
	task := &model.Task{}
	id := ctx.Param(ID)
	db := h.DB.Preload("Report")
	result := db.First(task, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Task{}
	r.With(task)

	ctx.JSON(http.StatusOK, r)
}

// List godoc
// @summary List all tasks.
// @description List all tasks.
// @tags get
// @produce json
// @success 200 {object} []api.Task
// @router /tasks [get]
func (h TaskHandler) List(ctx *gin.Context) {
	var list []model.Task
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = db.Preload("Report")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Task{}
	for i := range list {
		r := Task{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	ctx.JSON(http.StatusOK, resources)
}

// Create godoc
// @summary Create a task.
// @description Create a task.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Task
// @router /tasks [post]
// @param task body api.Task true "Task data"
func (h TaskHandler) Create(ctx *gin.Context) {
	task := Task{}
	err := ctx.BindJSON(&task)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}

	m := task.Model()
	m.Reset()
	result := h.DB.Create(&m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	task.With(m)

	ctx.JSON(http.StatusCreated, task)
}

// AddonCreate godoc
// @summary Create an addon task.
// @description Create an addon task.
// @tags create
// @accept json
// @produce json
// @success 201 {object} api.Task
// @router /addons/:name/tasks [post]
// @param task body api.Task true "Task data"
func (h TaskHandler) AddonCreate(ctx *gin.Context) {
	name := ctx.Param(Name)
	addon := &crd.Addon{}
	err := h.Client.Get(
		context.TODO(),
		client.ObjectKey{
			Namespace: Settings.Hub.Namespace,
			Name:      name,
		},
		addon)
	if err != nil {
		if errors.IsNotFound(err) {
			ctx.Status(http.StatusNotFound)
			return
		}
	}
	task := Task{}
	task.Name = addon.Name
	task.Addon = addon.Name
	task.Image = addon.Spec.Image
	err = ctx.BindJSON(&task.Data)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	m := task.Model()
	result := h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	task.With(m)

	ctx.JSON(http.StatusCreated, task)
}

// Delete godoc
// @summary Delete a task.
// @description Delete a task.
// @tags delete
// @success 204
// @router /tasks/:id [delete]
// @param id path string true "Task ID"
func (h TaskHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	task := &model.Task{}
	result := h.DB.First(task, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}
	if task.Job != "" {
		job := &batch.Job{}
		job.Namespace = path.Dir(task.Job)
		job.Name = path.Base(task.Job)
		err := h.Client.Delete(
			context.TODO(),
			job)
		if err != nil {
			h.deleteFailed(ctx, result.Error)
		}
	}
	result = h.DB.Delete(task, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a task.
// @description Update a task.
// @tags update
// @accept json
// @success 204
// @router /tasks/:id [put]
// @param id path string true "Task ID"
// @param task body Task true "Task data"
func (h TaskHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := &Task{}
	err := ctx.BindJSON(updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	m := updates.Model()
	result := h.DB.Model(&Task{}).Where("id", id).Omit("id").Updates(m)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// CreateReport godoc
// @summary Create a task report.
// @description Update a task report.
// @tags update
// @accept json
// @produce json
// @success 201 {object} api.TaskReport
// @router /tasks/:id [put]
// @param id path string true "TaskReport ID"
// @param task body api.TaskReport true "TaskReport data"
func (h TaskHandler) CreateReport(ctx *gin.Context) {
	id := ctx.Param(ID)
	report := &TaskReport{}
	err := ctx.BindJSON(report)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	task, _ := strconv.Atoi(id)
	report.TaskID = uint(task)
	m := report.Model()
	result := h.DB.Create(m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
	}
	report.With(m)

	ctx.JSON(http.StatusCreated, report)
}

// UpdateReport godoc
// @summary Update a task report.
// @description Update a task report.
// @tags update
// @accept json
// @produce json
// @success 200 {object} api.TaskReport
// @router /tasks/:id [put]
// @param id path string true "TaskReport ID"
// @param task body api.TaskReport true "TaskReport data"
func (h TaskHandler) UpdateReport(ctx *gin.Context) {
	id := ctx.Param(ID)
	report := &TaskReport{}
	err := ctx.BindJSON(report)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	task, _ := strconv.Atoi(id)
	report.TaskID = uint(task)
	m := report.Model()
	db := h.DB.Model(&model.TaskReport{})
	db = db.Where("taskid", task)
	db = db.Omit("id")
	result := db.Updates(m)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
	}
	report.With(m)

	ctx.JSON(http.StatusOK, report)
}

//
// Task REST resource.
type Task struct {
	ID         uint        `json:"id"`
	Name       string      `json:"name"`
	Image      string      `json:"image"`
	Addon      string      `json:"addon"`
	Data       model.JSON  `json:"data"`
	Started    *time.Time  `json:"started"`
	Terminated *time.Time  `json:"terminated"`
	Status     string      `json:"status"`
	Error      string      `json:"error"`
	Job        string      `json:"job"`
	Report     *TaskReport `json:"report"`
}

//
// With updates the resource with the model.
func (r *Task) With(m *model.Task) {
	r.ID = m.ID
	r.Name = m.Name
	r.Image = m.Image
	r.Addon = m.Addon
	r.Data = m.Data
	r.Started = m.Started
	r.Terminated = m.Terminated
	r.Status = m.Status
	r.Error = m.Error
	r.Job = m.Job
	if m.Report != nil {
		report := &TaskReport{}
		report.With(m.Report)
		r.Report = report
	}
}

//
// Model builds a model.
func (r *Task) Model() (m *model.Task) {
	m = &model.Task{
		Name:       r.Name,
		Image:      r.Image,
		Addon:      r.Addon,
		Data:       r.Data,
		Started:    r.Started,
		Terminated: r.Terminated,
		Status:     r.Status,
		Error:      r.Error,
		Job:        r.Job,
	}
	m.ID = r.ID
	if r.Report != nil {
		m.Report = r.Report.Model()
	}

	return
}

//
// TaskReport REST resource.
type TaskReport struct {
	ID        uint   `json:"id"`
	Error     string `json:"error"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
	Activity  string `json:"activity"`
	TaskID    uint   `json:"task"`
}

//
// With updates the resource with the model.
func (r *TaskReport) With(m *model.TaskReport) {
	r.ID = m.ID
	r.Error = m.Error
	r.Total = m.Total
	r.Completed = m.Completed
	r.Activity = m.Activity
	r.TaskID = m.TaskID
}

//
// Model builds a model.
func (r *TaskReport) Model() (m *model.TaskReport) {
	m = &model.TaskReport{
		Error:     r.Error,
		Total:     r.Total,
		Completed: r.Completed,
		Activity:  r.Activity,
		TaskID:    r.TaskID,
	}
	m.ID = r.ID

	return
}
