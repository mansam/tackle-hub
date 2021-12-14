package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/konveyor/tackle-hub/model"
	"mime"
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
	e.GET(AppBucketsRoot, h.List)
	e.GET(AppBucketsRoot+"/", h.List)
	e.POST(BucketsRoot, h.Create)
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
	bucket := Bucket{}
	id := ctx.Param(ID)
	result := h.DB.First(&bucket, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusOK, bucket)
}

// List godoc
// @summary List all buckets.
// @description List all buckets.
// @buckets get
// @produce json
// @success 200 {object} Bucket
// @router /controls/bucket [get]
func (h BucketHandler) List(ctx *gin.Context) {
	appId := ctx.Param(ID)
	var list []Bucket
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
// @summary Create a bucket.
// @description Create a bucket.
// @buckets create
// @accept json
// @produce json
// @success 200 {object} Bucket
// @router /controls/bucket [post]
// @param bucket body Bucket true "Bucket data"
func (h BucketHandler) Create(ctx *gin.Context) {
	bucket := Bucket{}
	err := ctx.BindJSON(&bucket)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	uid := uuid.New()
	bucket.Path = pathlib.Join(
		Settings.Hub.Bucket.Path,
		uid.String())
	err = os.MkdirAll(bucket.Path, 0755)
	if err != nil {
		h.createFailed(ctx, err)
		return
	}
	result := h.DB.Create(&bucket)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	ctx.JSON(http.StatusCreated, bucket)
}

// Delete godoc
// @summary Delete a bucket.
// @description Delete a bucket.
// @buckets delete
// @success 200 {object} Bucket
// @router /controls/bucket/:id [delete]
// @param id path string true "Bucket ID"
func (h BucketHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	bucket := &Bucket{}
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

	ctx.Status(http.StatusOK)
}

// Update godoc
// @summary Update a bucket.
// @description Update a bucket.
// @buckets update
// @accept json
// @produce json
// @success 200 {object} Bucket
// @router /controls/bucket/:id [put]
// @param id path string true "Bucket ID"
// @param bucket body Bucket true "Bucket data"
func (h BucketHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)
	updates := Bucket{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		h.updateFailed(ctx, err)
		return
	}
	db := h.DB.Model(&Bucket{})
	db = db.Where("id", id)
	db = db.Omit("id", "location")
	result := db.Updates(updates)
	if result.Error != nil {
		h.updateFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h BucketHandler) GetContent(ctx *gin.Context) {
	path := ctx.Param(Wildcard)
	bucket := Bucket{}
	id := ctx.Param(ID)
	result := h.DB.First(&bucket, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	b, err := os.ReadFile(pathlib.Join(
		bucket.Path,
		path))
	if err != nil {
		h.getFailed(ctx, err)
		return
	}
	ctx.Data(
		http.StatusOK,
		mime.TypeByExtension(pathlib.Ext(path)),
		b)
}

//
// Bucket REST resource.
type Bucket = model.Bucket
