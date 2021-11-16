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
	BinaryReposRoot = InventoryRoot + "/binary-repository"
	BinaryRepoRoot  = BinaryReposRoot + "/:" + ID
)

type BinaryRepoHandler struct{}

func (h *BinaryRepoHandler) AddRoutes(e *gin.Engine) {
	e.GET(BinaryReposRoot, h.List)
	e.GET(BinaryReposRoot+"/", h.List)
	e.POST(BinaryReposRoot, h.Create)
	e.GET(BinaryRepoRoot, h.Get)
	e.PUT(BinaryRepoRoot, h.Update)
	e.DELETE(BinaryRepoRoot, h.Delete)
}

func (h *BinaryRepoHandler) Get(ctx *gin.Context) {
	binaryRepo := models.BinaryRepo{}
	id := ctx.Param(ID)
	result := db.DB.First(&binaryRepo, "id = ?", id)
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
	ctx.JSON(http.StatusOK, binaryRepo)
}

func (h *BinaryRepoHandler) List(ctx *gin.Context) {
	var binaryRepos []models.BinaryRepo
	result := db.DB.Find(&binaryRepos)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, binaryRepos)
}

func (h *BinaryRepoHandler) Create(ctx *gin.Context) {
	binaryRepo := models.BinaryRepo{}
	err := ctx.BindJSON(&binaryRepo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := db.DB.Create(&binaryRepo)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": MsgInternalServerError,
		})
		log.Error(result.Error, MsgInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, binaryRepo)
}

func (h *BinaryRepoHandler) Delete(ctx *gin.Context) {
	id := ctx.Param(ID)

	result := db.DB.Delete(&models.BinaryRepo{}, "id = ?", id)
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

func (h *BinaryRepoHandler) Update(ctx *gin.Context) {
	id := ctx.Param(ID)

	updates := models.BinaryRepo{}
	err := ctx.BindJSON(&updates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": MsgBadRequest,
		})
		log.Error(err, MsgBadRequest)
		return
	}

	result := db.DB.Model(&models.BinaryRepo{}).Where("id = ?", id).Omit("id").Updates(updates)
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
