package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
	"os"
	pathlib "path"
)

//
// Routes
const (
	BucketsRoot    = "/buckets"
	BucketRoot     = BucketsRoot + "/:" + ID
	BucketContent  = BucketRoot + "/*" + Wildcard
	AppBucketsRoot = ApplicationRoot + BucketsRoot
)

//
// BucketHandler handles bucket routes.
type BucketHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h BucketHandler) AddRoutes(e *gin.Engine) {
	e.GET(BucketsRoot, h.List)
	e.GET(BucketsRoot+"/", h.List)
	e.GET(AppBucketsRoot, h.ListByApplication)
	e.GET(AppBucketsRoot+"/", h.ListByApplication)
	e.POST(BucketsRoot, h.Create)
	e.POST(AppBucketsRoot, h.CreateForApplication)
	e.GET(BucketRoot, h.Get)
	e.PUT(BucketRoot, h.Update)
	e.DELETE(BucketRoot, h.Delete)
	e.GET(BucketContent, h.GetContent)
}

// Get godoc
// @summary Get a bucket by ID.
// @description Get a bucket by ID.
// @buckets get
// @produce json
// @success 200 {object} Bucket
// @router /controls/bucket/:id [get]
// @param id path string true "Bucket ID"
func (h BucketHandler) Get(ctx *gin.Context) {
	m := &model.Bucket{}
	id := ctx.Param(ID)
	result := h.DB.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	r := Bucket{}
	r.With(m)

	ctx.JSON(http.StatusOK, r)
}

// List godoc
// @summary List all buckets.
// @description List all buckets.
// @buckets get
// @produce json
// @success 200 {object} []Bucket
// @router /controls/bucket [get]
func (h BucketHandler) List(ctx *gin.Context) {
	var list []model.Bucket
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Bucket{}
	for i := range list {
		r := Bucket{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	ctx.JSON(http.StatusOK, resources)
}

// ListByApplication godoc
// @summary List all buckets.
// @description List all buckets.
// @buckets get
// @produce json
// @success 200 {object} []Bucket
// @router /controls/bucket [get]
func (h BucketHandler) ListByApplication(ctx *gin.Context) {
	var list []model.Bucket
	appId := ctx.Param(ID)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = db.Where("applicationid", appId)
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Bucket{}
	for i := range list {
		r := Bucket{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	ctx.JSON(http.StatusOK, resources)
}

// Create godoc
// @summary Create a bucket.
// @description Create a bucket.
// @buckets create
// @accept json
// @produce json
// @success 201 {object} Bucket
// @router /controls/bucket [post]
// @param bucket body Bucket true "Bucket data"
func (h BucketHandler) Create(ctx *gin.Context) {
	bucket := &Bucket{}
	err := ctx.BindJSON(bucket)
	if err != nil {
		return
	}
	err = h.create(bucket)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, bucket)
}

// CreateForApplication godoc
// @summary Create a bucket for an application.
// @description Create a bucket for an application.
// @buckets create
// @accept json
// @produce json
// @success 201 {object} Bucket
// @router /controls/bucket [post]
// @param bucket body Bucket true "Bucket data"
func (h BucketHandler) CreateForApplication(ctx *gin.Context) {
	bucket := &Bucket{}
	err := ctx.BindJSON(bucket)
	if err != nil {
		return
	}
	appID := ctx.Param(ID)
	application := &model.Application{}
	result := h.DB.First(application, appID)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	bucket.ApplicationID = application.ID
	err = h.create(bucket)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, bucket)
}

// Delete godoc
// @summary Delete a bucket.
// @description Delete a bucket.
// @buckets delete
// @success 204 {object} Bucket
// @router /controls/bucket/:id [delete]
// @param id path string true "Bucket ID"
func (h BucketHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	bucket := &model.Bucket{}
	result := h.DB.First(bucket, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}
	result = h.DB.Delete(bucket, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Update godoc
// @summary Update a bucket.
// @description Update a bucket.
// @buckets update
// @accept json
// @produce json
// @success 204 {object} Bucket
// @router /controls/bucket/:id [put]
// @param id path string true "Bucket ID"
// @param bucket body Bucket true "Bucket data"
func (h BucketHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	r := &Bucket{}
	err := ctx.BindJSON(r)
	if err != nil {
		h.bindFailed(ctx, err)
		return
	}
	m := r.Model()
	db := h.DB.Model(&model.Bucket{})
	db = db.Where("id", id)
	db = db.Omit("id", "location")
	result := db.Updates(m)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// GetContent streams file file content.
func (h BucketHandler) GetContent(ctx *gin.Context) {
	rPath := ctx.Param(Wildcard)
	bucket := model.Bucket{}
	id := ctx.Param(ID)
	result := h.DB.First(&bucket, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.File(pathlib.Join(
		bucket.Path,
		rPath))
}

//
// create a bucket.
func (h BucketHandler) create(bucket *Bucket) (err error) {
	uid := uuid.New()
	bucket.Path = pathlib.Join(
		Settings.Hub.Bucket.Path,
		uid.String())
	err = os.MkdirAll(bucket.Path, 0777)
	if err != nil {
		return
	}

	m := bucket.Model()
	result := h.DB.Create(m)
	err = result.Error
	if err != nil {
		_ = os.Remove(bucket.Path)
	}
	bucket.With(m)

	return
}

//
// Bucket REST Resource.
type Bucket struct {
	Resource
	Name          string `json:"name"`
	Path          string `json:"path"`
	ApplicationID uint   `json:"application"`
}

//
// With updates the resource with the model.
func (r *Bucket) With(m *model.Bucket) {
	r.Resource.With(&m.Model)
	r.Name = m.Name
	r.Path = m.Path
	r.ApplicationID = m.ApplicationID
}

//
// Model builds a model.
func (r *Bucket) Model() (m *model.Bucket) {
	m = &model.Bucket{
		Name:          r.Name,
		Path:          r.Path,
		ApplicationID: r.ApplicationID,
	}
	m.ID = r.ID

	return
}
