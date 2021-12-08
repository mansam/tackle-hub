package addon

import (
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
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
func (h *Application) Get(id uint) (m *model.Application, err error) {
	m = &model.Application{}
	err = h.client.Get(
		pathlib.Join(
			api.ApplicationsRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// Update an application by ID.
func (h *Application) Update(m *model.Application) (err error) {
	m = &model.Application{}
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
	artifact := &model.Artifact{}
	artifact.CreateUser = "addon"
	artifact.Name = pathlib.Base(path)
	artifact.ApplicationID = application
	artifact.Kind = kind
	err = h.client.Post(api.ArtifactsRoot, artifact)
	return
}

//
// Tag API.
type Tag struct {
	// hub API client.
	client *Client
}

//
// Create a tag.
func (h *Tag) Create(tagType uint, name string) (tag *model.Tag, err error) {
	tag = &model.Tag{}
	tag.TagTypeID = tagType
	tag.Name = name
	err = h.client.Post(api.TagsRoot, tag)
	return
}
