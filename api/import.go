package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
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
// Import Statuses
const (
	InProgress = "In Progress"
	Completed  = "Completed"
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
	SummariesRoot = InventoryRoot + "/import-summary"
	SummaryRoot   = SummariesRoot + "/:" + ID
	ImportsRoot   = InventoryRoot + "/application-import"
	ImportRoot    = ImportsRoot + "/:" + ID
	DownloadRoot  = InventoryRoot + "/csv-export"
)

//
// ImportHandler handles import routes.
type ImportHandler struct {
	BaseHandler
}

//
// AddRoutes adds routes.
func (h ImportHandler) AddRoutes(e *gin.Engine) {
	e.GET(SummariesRoot, h.ListSummaries)
	e.GET(SummariesRoot+"/", h.ListSummaries)
	e.GET(SummaryRoot, h.GetSummary)
	e.DELETE(SummaryRoot, h.DeleteSummary)
	e.GET(ImportsRoot, h.ListImports)
	e.GET(ImportsRoot+"/", h.ListImports)
	e.GET(ImportRoot, h.GetImport)
	e.DELETE(ImportRoot, h.DeleteImport)
	e.GET(DownloadRoot, h.DownloadCSV)
	e.POST(UploadRoot, h.UploadCSV)
}

//
// GetImport godoc
// @summary Get an import by ID.
// @description Get an import by ID.
// @tags get
// @produce json
// @success 200 {object} api.Import
// @router /application-inventory/application-import/{id} [get]
// @param id path string true "Import ID"
func (h ImportHandler) GetImport(ctx *gin.Context) {
	m := &model.Import{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "ImportTags")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, m.AsMap())
}

//
// ListImports godoc
// @summary List imports.
// @description List imports.
// @tags list
// @produce json
// @success 200 {object} []api.Import
// @router /application-inventory/application-import [get]
func (h ImportHandler) ListImports(ctx *gin.Context) {
	var count int64
	var list []model.Import
	db := h.DB
	summaryId := ctx.Query("importSummary.id")
	if summaryId != "" {
		db = db.Where("importsummaryid = ?", summaryId)
	}
	isValid := ctx.Query("isValid")
	if isValid == "true" {
		db = db.Where("isvalid")
	} else if isValid == "false" {
		db = db.Not("isvalid")
	}
	db.Model(model.Import{}).Count(&count)
	pagination := NewPagination(ctx)
	db = pagination.apply(db)
	db = h.preLoad(db, "ImportTags")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []Import{}
	for i := range list {
		resources = append(resources, list[i].AsMap())
	}

	h.listResponse(ctx, ImportKind, resources, int(count))
}

//
// DeleteImport godoc
// @summary Delete an import.
// @description Delete an import. This leaves any created application or dependency.
// @tags delete
// @success 204
// @router /application-inventory/application-import/{id} [delete]
// @param id path string true "Import ID"
func (h ImportHandler) DeleteImport(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.Import{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// GetSummary godoc
// @summary Get an import summary by ID.
// @description Get an import by ID.
// @tags get
// @produce json
// @success 200 {object} api.ImportSummary
// @router /application-inventory/import-summary/{id} [get]
// @param id path string true "ImportSummary ID"
func (h ImportHandler) GetSummary(ctx *gin.Context) {
	m := &model.ImportSummary{}
	id := ctx.Param(ID)
	db := h.preLoad(h.DB, "Imports")
	result := db.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.JSON(http.StatusOK, m)
}

//
// ListSummaries godoc
// @summary List import summaries.
// @description List import summaries.
// @tags list
// @produce json
// @success 200 {object} []api.ImportSummary
// @router /application-inventory/import-summary [get]
func (h ImportHandler) ListSummaries(ctx *gin.Context) {
	var count int64
	var list []model.ImportSummary
	h.DB.Model(model.ImportSummary{}).Count(&count)
	pagination := NewPagination(ctx)
	db := pagination.apply(h.DB)
	db = h.preLoad(db, "Imports")
	result := db.Find(&list)
	if result.Error != nil {
		h.listFailed(ctx, result.Error)
		return
	}
	resources := []ImportSummary{}
	for i := range list {
		r := ImportSummary{}
		r.With(&list[i])
		resources = append(resources, r)
	}

	h.listResponse(ctx, SummaryKind, resources, int(count))
}

//
// DeleteSummary godoc
// @summary Delete an import summary and associated import records.
// @description Delete an import summary and associated import records.
// @tags delete
// @success 204
// @router /application-inventory/import-summary/{id} [delete]
// @param id path string true "ImportSummary ID"
func (h ImportHandler) DeleteSummary(ctx *gin.Context) {
	id := ctx.Param(ID)
	result := h.DB.Delete(&model.ImportSummary{}, id)
	if result.Error != nil {
		h.deleteFailed(ctx, result.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

//
// UploadCSV godoc
// @summary Upload a CSV containing applications and dependencies to import.
// @description Upload a CSV containing applications and dependencies to import.
// @tags post
// @success 201 {object} api.ImportSummary
// @produce json
// @router /application-inventory/file/upload [post]
func (h ImportHandler) UploadCSV(ctx *gin.Context) {
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
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, fileReader)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	m := model.ImportSummary{
		Filename:     fileName,
		ImportStatus: InProgress,
		Content:      buf.Bytes(),
	}
	result := h.DB.Create(&m)
	if result.Error != nil {
		h.createFailed(ctx, result.Error)
		return
	}
	_, err = fileReader.Seek(0, 0)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	csvReader := csv.NewReader(fileReader)
	csvReader.TrimLeadingSpace = true
	// skip the header
	_, err = csvReader.Read()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
	}

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ctx.Status(http.StatusBadRequest)
			}
		}
		var imp model.Import
		switch row[0] {
		case RecordTypeApplication:
			imp = h.applicationFromRow(fileName, row)
		case RecordTypeDependency:
			imp = h.dependencyFromRow(fileName, row)
		default:
			imp = model.Import{
				Filename:    fileName,
				RecordType1: row[0],
			}
		}
		imp.ImportSummary = m
		result := h.DB.Create(&imp)
		if result.Error != nil {
			h.createFailed(ctx, result.Error)
			return
		}
	}

	summary := ImportSummary{}
	summary.With(&m)
	ctx.JSON(http.StatusCreated, summary)
}

//
// DownloadCSV godoc
// @summary Export the source CSV for a particular import summary.
// @description Export the source CSV for a particular import summary.
// @tags export
// @produce text/csv
// @success 200 file csv
// @router /application-inventory/csv-export [get]
// @param importSummary.id query string true "ImportSummary ID"
func (h ImportHandler) DownloadCSV(ctx *gin.Context) {
	id := ctx.Query("importSummary.id")
	m := &model.ImportSummary{}
	result := h.DB.First(m, id)
	if result.Error != nil {
		h.getFailed(ctx, result.Error)
		return
	}
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", m.Filename))
	ctx.Data(http.StatusOK, "text/csv", m.Content)
}

//
// CSV upload supports two types of records in the same file: application imports, and dependencies.
// A dependency row must consist of the following columns:
//
// Col 1: Record Type 1 -- This will always contain a "2" for a dependency
// Col 2: Application Name -- The name of the application that has the dependency relationship.
//                            This application must exist.
// Col N-2: Dependency -- The name of the application on the other side of the dependency relationship.
// Col N-1: Dependency Direction -- Whether this is a "northbound" or "southbound" dependency.
//
// Between the Application Name and the Dependency field there may be an arbitrary number of columns representing
// tags or other fields that only pertain to an application import. The dependency and direction will always be
// the last two columns in the row.
//
// Examples:
//
// 2,MyApplication,,,,,,,,,,,OtherApplication,NORTHBOUND
// 2,MyApplication,OtherApplication,SOUTHBOUND
func (h ImportHandler) dependencyFromRow(fileName string, row []string) (app model.Import) {
	app = model.Import{
		Filename:            fileName,
		RecordType1:         row[0],
		ApplicationName:     row[1],
		Dependency:          row[len(row)-2],
		DependencyDirection: row[len(row)-1],
	}
	return
}

//
// CSV upload supports two types of records in the same file: application imports, and dependencies.
// An application row must consist of the following columns:
//
// Col 1: Record Type 1 -- This will always contain a "1" for an application
// Col 2: Application Name -- The name of the application to be created.
// Col 3: Description -- A short description of the application.
// Col 4: Comments -- Additional comments on the application.
// Col 5: Business Service -- The name of the business service this Application should belong to.
//                            This business service must already exist.
//
// Following that are up to twenty pairs of Tag Types and Tags, specified by name. These are optional.
// If a tag type and a tag are specified, they must already exist.
//
// Examples:
//
// 1,MyApplication,My cool app,No comment,MyBusinessService,TagType1,Tag1,TagType2,Tag2
// 1,OtherApplication,,,MyBusinessService
func (h ImportHandler) applicationFromRow(fileName string, row []string) (app model.Import) {
	app = model.Import{
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
				Name:    row[i],
				TagType: row[i-1],
			}
			app.ImportTags = append(app.ImportTags, tag)
		}
	}

	return
}

//
// Import REST resource.
type Import map[string]interface{}

//
// ImportSummary REST resource.
type ImportSummary struct {
	Resource
	Filename     string    `json:"filename"`
	ImportStatus string    `json:"importStatus"`
	ImportTime   time.Time `json:"importTime"`
	ValidCount   int       `json:"validCount"`
	InvalidCount int       `json:"invalidCount"`
}

//
// With updates the resource with the model.
func (r *ImportSummary) With(m *model.ImportSummary) {
	r.Resource.With(&m.Model)
	r.Filename = m.Filename
	r.ImportTime = m.CreateTime
	for _, imp := range m.Imports {
		if imp.Processed {
			if imp.IsValid {
				r.ValidCount++
			} else {
				r.InvalidCount++
			}
		}
	}
	if len(m.Imports) == r.ValidCount+r.InvalidCount {
		r.ImportStatus = Completed
	} else {
		r.ImportStatus = InProgress
	}
}
