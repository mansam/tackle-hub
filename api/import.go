package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"net/http"
)

//
// Kind
const (
	ImportKind = "application-import"
)

//
// Routes
const (
	UploadRoot  = InventoryRoot + "/file/upload"
	ImportsRoot = InventoryRoot + "/application-import"
	ImportRoot  = ImportsRoot + "/:" + ID
)

type ImportHandler struct {
	BaseHandler
}

func (h ImportHandler) AddRoutes(e *gin.Engine) {
	e.GET(ImportsRoot, h.List)
	e.GET(ImportsRoot+"/", h.List)
	e.POST(UploadRoot, h.Create)
	e.GET(ImportRoot, h.Get)
	e.DELETE(ImportRoot, h.Delete)
}

func (h ImportHandler) Get(ctx *gin.Context) {
	m := &model.ApplicationImport{}
	id := ctx.Param(ID)
	result := h.DB.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, m)
}

func (h ImportHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.ApplicationImport
	h.DB.Model(model.ApplicationImport{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}

	list := List{}
	list.With(ImportKind, models, int(count))
	h.hal(ctx, http.StatusOK, list)
}

func (h ImportHandler) Create(ctx *gin.Context) {
	//fileName, ok := ctx.GetPostForm("fileName")
	//if !ok {
	//	ctx.Status(http.StatusBadRequest)
	//}
	//file, ok := ctx.GetPostForm("file")
	//if !ok {
	//	ctx.Status(http.StatusBadRequest)
	//}
	//
	//imports := []model.ApplicationImport{}
	//
	//reader := csv.NewReader(strings.NewReader(file))
	//// skip the header
	//_, err := reader.Read()
	//if err != nil {
	//	ctx.Status(http.StatusBadRequest)
	//}
	//for {
	//	row, err := reader.Read()
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		} else {
	//			ctx.Status(http.StatusBadRequest)
	//		}
	//	}
	//	appImport := []model.ApplicationImport{}
	//
	//}

}

func (h ImportHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.ApplicationImport{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusOK)
}
