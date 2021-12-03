package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/konveyor/tackle-hub/k8s"
	crd "github.com/konveyor/tackle-hub/k8s/api/tackle/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//
// Routes
const (
	AddonsRoot = "/addons"
	AddonRoot  = AddonsRoot + "/:" + Name
)

//
// AddonHandler handles addon routes.
type AddonHandler struct {
	BaseHandler
	Client client.Client
}

//
// AddRoutes adds routes.
func (h AddonHandler) AddRoutes(e *gin.Engine) {
	e.GET(AddonsRoot, h.List)
	e.GET(AddonsRoot+"/", h.List)
	e.GET(AddonRoot, h.Get)
}

func (h AddonHandler) Get(ctx *gin.Context) {
	name := ctx.Param(Name)
	addon := &crd.Addon{}
	err := h.Client.Get(
		context.TODO(),
		client.ObjectKey{
			Namespace: k8s.Namespace,
			Name:      name,
		},
		addon)
	if err != nil {
		if errors.IsNotFound(err) {
			ctx.Status(http.StatusNotFound)
			return
		} else {
			h.getFailed(ctx, err)
			return
		}
	}
	r := Addon{}
	r.With(addon)

	ctx.JSON(http.StatusOK, r)
}

func (h AddonHandler) List(ctx *gin.Context) {
	list := &crd.AddonList{}
	err := h.Client.List(
		context.TODO(),
		nil,
		list)
	if err != nil {
		h.listFailed(ctx, err)
		return
	}
	content := []Addon{}
	for _, m := range list.Items {
		addon := Addon{}
		addon.With(&m)
		content = append(content, addon)
	}

	ctx.JSON(http.StatusOK, content)
}

//
// Addon REST resource.
type Addon struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

//
// With model.
func (r *Addon) With(m *crd.Addon) {
	r.Name = m.Name
	r.Image = m.Spec.Image
}
