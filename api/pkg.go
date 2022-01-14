package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/gorm"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

var (
	Settings = &settings.Settings
	log      = logging.WithName("api")
)

//
// Routes
const (
	InventoryRoot = "/application-inventory"
	ControlsRoot  = "/controls"
)

//
// Params
const (
	ID       = "id"
	Name     = "name"
	Wildcard = "Wildcard"
)

//
// Pagination Defaults
const (
	Limit  = 20
	Offset = 0
	Sort   = "id asc"
)

//
// All builds all handlers.
func All() []Handler {
	return []Handler{
		&ApplicationHandler{},
		&BucketHandler{},
		&BusinessServiceHandler{},
		&DependencyHandler{},
		&ImportHandler{},
		&JobFunctionHandler{},
		&RepositoryHandler{},
		&IdentityHandler{},
		&ReviewHandler{},
		&StakeholderHandler{},
		&StakeholderGroupHandler{},
		&TagHandler{},
		&TagTypeHandler{},
		&TaskHandler{},
		&AddonHandler{},
	}
}

//
// Handler.
type Handler interface {
	With(*gorm.DB, client.Client)
	AddRoutes(e *gin.Engine)
}

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

//
// Hal REST resource.
type Hal struct {
	Embedded   map[string]interface{} `json:"_embedded"`
	TotalCount int                    `json:"total_count"`
}

//
// With sets the embedded resource and count.
func (r *Hal) With(kind string, resources interface{}, total int) {
	r.Embedded = make(map[string]interface{})
	r.Embedded[kind] = resources
	r.TotalCount = total
}
