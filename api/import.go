package api

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/model"
	"io"
	"net/http"
	"time"
)

//
// Record types
const (
	RecordTypeApplication = "1"
	RecordTypeDependency  = "2"
)

//
// Kind
const (
	ImportKind  = "application-import"
	SummaryKind = "import-summary"
)

//
// Routes
const (
	UploadRoot    = InventoryRoot + "/file/upload"
	ExportRoot    = InventoryRoot + "/csv-export"
	SummariesRoot = InventoryRoot + "/import-summary"
	SummaryRoot   = SummariesRoot + "/:" + ID
	ImportsRoot   = InventoryRoot + "/application-import"
	ImportRoot    = ImportsRoot + "/:" + ID
)

type UploadHandler struct {
	BaseHandler
}

func (h UploadHandler) AddRoutes(e *gin.Engine) {
	e.POST(UploadRoot, h.Create)
}

func (h UploadHandler) Create(ctx *gin.Context) {
	fileName, ok := ctx.GetPostForm("fileName")
	if !ok {
		ctx.Status(http.StatusBadRequest)
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}
	fileReader, err := file.Open()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}
	reader := csv.NewReader(fileReader)
	// skip the header
	_, err = reader.Read()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}

	m := model.ImportSummary{
		Filename:     fileName,
		ImportStatus: "In Progress",
	}
	result := h.DB.Create(&m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ctx.Status(http.StatusBadRequest)
			}
		}
		var imp *model.ApplicationImport
		switch row[0] {
		case RecordTypeApplication:
			imp = h.importFromRow(fileName, row)
		case RecordTypeDependency:
			imp = h.dependencyFromRow(fileName, row)
		default:
			imp = &model.ApplicationImport{
				Filename:    fileName,
				RecordType1: row[0],
			}
		}
		imp.ImportSummary = m
		result := h.DB.Create(imp)
		if result.Error != nil {
			h.createFailed(ctx, result.Error)
			return
		}
	}

	summary := ImportSummary{}
	summary.With(&m)
	ctx.JSON(http.StatusCreated, summary)
}

func (h *UploadHandler) dependencyFromRow(fileName string, row []string) (app *model.ApplicationImport) {
	app = &model.ApplicationImport{
		Filename:            fileName,
		RecordType1:         row[0],
		Dependency:          row[len(row)-2],
		DependencyDirection: row[len(row)-1],
	}
	return
}

func (h *UploadHandler) importFromRow(fileName string, row []string) (app *model.ApplicationImport) {
	app = &model.ApplicationImport{
		Filename:        fileName,
		RecordType1:     row[0],
		ApplicationName: row[1],
		Description:     row[2],
		Comments:        row[3],
		BusinessService: row[4],
	}

	for i := 5; i < len(row); i++ {
		if i%2 == 0 {
			tag := model.ImportTag{
				Name:    row[i-1],
				TagType: row[i],
				Order:   uint(i - 5),
			}
			app.ImportTags = append(app.ImportTags, tag)
		}
	}

	return
}

type SummaryHandler struct {
	BaseHandler
}

func (h SummaryHandler) AddRoutes(e *gin.Engine) {
	e.GET(SummariesRoot, h.List)
	e.GET(SummariesRoot+"/", h.List)
	e.GET(SummaryRoot, h.Get)
	e.DELETE(SummaryRoot, h.Delete)
}

func (h SummaryHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.ImportSummary
	h.DB.Model(model.ImportSummary{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "ApplicationImports")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []ImportSummary{}
	for i := range models {
		r := ImportSummary{}
		r.With(&models[i])
		resources = append(resources, r)
	}

	list := List{}
	list.With(SummaryKind, resources, int(count))
	h.hal(ctx, http.StatusOK, list)
}

func (h SummaryHandler) Get(ctx *gin.Context) {
	m := &model.ImportSummary{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "ApplicationImports")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, m)
}

func (h SummaryHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.ImportSummary{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

type ImportHandler struct {
	BaseHandler
}

func (h ImportHandler) AddRoutes(e *gin.Engine) {
	e.GET(ImportsRoot, h.List)
	e.GET(ImportsRoot+"/", h.List)
	e.GET(ImportRoot, h.Get)
	e.DELETE(ImportRoot, h.Delete)
}

func (h ImportHandler) Get(ctx *gin.Context) {
	m := &model.ApplicationImport{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "ImportTags")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, m.AsMap())
}

func (h ImportHandler) List(ctx *gin.Context) {
	var count int64
	var models []model.ApplicationImport
	summaryId := ctx.Query("importSummary.id")
	db := h.DB.Where("import_summary_id = ? AND is_valid = false", summaryId)
	db.Model(model.ApplicationImport{}).Count(&count)
	pagination := NewPagination(ctx)
	db = pagination.apply(db)
	db = h.preLoad(db, "ImportTags")
	result := db.Find(&models)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []map[string]interface{}{}
	for i := range models {
		resources = append(resources, models[i].AsMap())
	}

	list := List{}
	list.With(ImportKind, resources, int(count))
	h.hal(ctx, http.StatusOK, list)
}

func (h ImportHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.ApplicationImport{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

type ExportHandler struct {
	BaseHandler
}

func (h ExportHandler) AddRoutes(e *gin.Engine) {
	e.GET(ExportRoot, h.Get)
}

func (h ExportHandler) Get(ctx *gin.Context) {
	summaryId := ctx.Query("importSummary.id")
	imports := &[]model.ApplicationImport{}
	result := h.DB.Find(imports).Where("import_summary_id = ?", summaryId)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	// TODO:
	ctx.Status(http.StatusNoContent)
}

type ImportSummary struct {
	model.ImportSummary
	ImportTime   time.Time `json:"importTime"`
	ValidCount   int       `json:"validCount"`
	InvalidCount int       `json:"invalidCount"`
}

func (r *ImportSummary) With(m *model.ImportSummary) {
	r.ImportSummary = *m
	r.ImportTime = m.CreateTime
	for _, imp := range r.ApplicationImports {
		if imp.Processed {
			if imp.IsValid {
				r.ValidCount++
			} else {
				r.InvalidCount++
			}
		}
	}
	if len(r.ApplicationImports) == r.InvalidCount {
		r.ImportStatus = "Failed"
	} else if len(r.ApplicationImports) == r.ValidCount+r.InvalidCount {
		r.ImportStatus = "Completed"
	} else {
		r.ImportStatus = "In Progress"
	}
}
