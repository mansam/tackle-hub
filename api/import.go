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
		switch i {
		case 6:
			app.Tag1 = row[i-1]
			app.TagType1 = row[i]
		case 8:
			app.Tag1 = row[i-1]
			app.TagType1 = row[i]
		case 10:
			app.Tag2 = row[i-1]
			app.TagType2 = row[i]
		case 12:
			app.Tag3 = row[i-1]
			app.TagType3 = row[i]
		case 14:
			app.Tag4 = row[i-1]
			app.TagType4 = row[i]
		case 16:
			app.Tag5 = row[i-1]
			app.TagType5 = row[i]
		case 18:
			app.Tag6 = row[i-1]
			app.TagType6 = row[i]
		case 20:
			app.Tag7 = row[i-1]
			app.TagType7 = row[i]
		case 22:
			app.Tag8 = row[i-1]
			app.TagType8 = row[i]
		case 24:
			app.Tag9 = row[i-1]
			app.TagType9 = row[i]
		case 26:
			app.Tag10 = row[i-1]
			app.TagType10 = row[i]
		case 28:
			app.Tag11 = row[i-1]
			app.TagType11 = row[i]
		case 30:
			app.Tag12 = row[i-1]
			app.TagType12 = row[i]
		case 32:
			app.Tag13 = row[i-1]
			app.TagType13 = row[i]
		case 34:
			app.Tag13 = row[i-1]
			app.TagType13 = row[i]
		case 36:
			app.Tag14 = row[i-1]
			app.TagType14 = row[i]
		case 38:
			app.Tag15 = row[i-1]
			app.TagType15 = row[i]
		case 40:
			app.Tag16 = row[i-1]
			app.TagType16 = row[i]
		case 42:
			app.Tag17 = row[i-1]
			app.TagType18 = row[i]
		case 44:
			app.Tag19 = row[i-1]
			app.TagType19 = row[i]
		case 46:
			app.Tag20 = row[i-1]
			app.TagType20 = row[i]
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
	result := h.DB.First(m, id)
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
	id := ctx.Query("importSummaryId")
	imports := &[]model.ApplicationImport{}
	result := h.DB.Find(imports).Where("import_summary_id = ?", id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	// TODO: find csv file on disk and return it.
	ctx.Status(http.StatusNoContent)
}

type ImportSummary struct {
	model.ImportSummary
	ImportTime time.Time `json:"importTime"`
}

func (r *ImportSummary) With(m *model.ImportSummary) {
	r.ImportSummary = *m
	r.ImportTime = m.CreateTime
}
