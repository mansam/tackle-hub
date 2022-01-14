package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

//
// Pagination Defaults
const (
	Limit  = 20
	Offset = 0
	Sort   = "id asc"
)

//
// Pagination provides pagination and sorting.
type Pagination struct {
	Limit  int
	Offset int
	Sort   string
}

//
// apply pagination.
func (p *Pagination) apply(db *gorm.DB) (tx *gorm.DB) {
	tx = db.Offset(p.Offset).Limit(p.Limit)
	tx = tx.Order(p.Sort)
	return
}

//
// NewPagination factory.
func NewPagination(ctx *gin.Context) Pagination {
	limit, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		limit = Limit
	}
	offset, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		offset = Offset
	}
	sort := ctx.Query("sort")
	if sort == "" {
		sort = Sort
	}
	return Pagination{
		Limit:  limit,
		Offset: offset * limit,
		Sort:   sort,
	}
}
