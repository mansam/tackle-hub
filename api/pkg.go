package api

import (
	"github.com/gin-gonic/gin"
	"github.com/konveyor/controller/pkg/logging"
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/gorm"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
// All builds all handlers.
func All() []Handler {
	return []Handler{
		&AddonHandler{},
		&ApplicationHandler{},
		&BucketHandler{},
		&BusinessServiceHandler{},
		&DependencyHandler{},
		&ImportHandler{},
		&JobFunctionHandler{},
		&RepositoryHandler{},
		&IdentityHandler{},
		&ProxyHandler{},
		&ReviewHandler{},
		&StakeholderHandler{},
		&StakeholderGroupHandler{},
		&TagHandler{},
		&TagTypeHandler{},
		&TaskHandler{},
	}
}

//
// Handler.
type Handler interface {
	With(*gorm.DB, client.Client)
	AddRoutes(e *gin.Engine)
}