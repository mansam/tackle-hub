package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
	"strconv"
)

//
// Application API.
type Application struct {
	// hub API client.
	client *Client
}

//
// Get an application by ID.
func (h *Application) Get(id uint) (m *api.Application, err error) {
	m = &api.Application{}
	err = h.client.Get(
		pathlib.Join(
			api.ApplicationsRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List applications.
func (h *Application) List() (list []api.Application, err error) {
	list = []api.Application{}
	err = h.client.Get(api.ApplicationsRoot, &list)
	return
}

//
// Update an application by ID.
func (h *Application) Update(m *api.Application) (err error) {
	m = &api.Application{}
	err = h.client.Put(
		pathlib.Join(
			api.ApplicationsRoot,
			strconv.Itoa(int(m.ID))),
		m)
	return
}

//
// Artifact API.
type Artifact struct {
	// hub API client.
	client *Client
}

//
// Upload an artifact.
func (h *Artifact) Upload(application uint, kind string, path string) (err error) {
	artifact := &api.Artifact{}
	artifact.CreateUser = "addon"
	artifact.Name = pathlib.Base(path)
	artifact.ApplicationID = application
	artifact.Kind = kind
	err = h.client.Post(api.ArtifactsRoot, artifact)
	return
}
