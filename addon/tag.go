package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
	"strconv"
)

//
// Tag API.
type Tag struct {
	// hub API client.
	client *Client
}

//
// Create a tag.
func (h *Tag) Create(m *api.Tag) (err error) {
	err = h.client.Post(api.TagsRoot, m)
	if err == nil {
		Log.Info(
			"Addon created: tag.",
			"object",
			m)
	}
	return
}

//
// Get a tag by ID.
func (h *Tag) Get(id uint) (m *api.Tag, err error) {
	m = &api.Tag{}
	err = h.client.Get(
		pathlib.Join(
			api.TagsRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List tags.
func (h *Tag) List() (list []api.Tag, err error) {
	list = []api.Tag{}
	err = h.client.Get(api.TagsRoot, &list)
	return
}

//
// Delete a tag.
func (h *Tag) Delete(m *api.Tag) (err error) {
	err = h.client.Delete(
		pathlib.Join(
			api.TagsRoot,
			strconv.Itoa(int(m.ID))))
	if err == nil {
		Log.Info(
			"Addon deleted: tag.",
			"object",
			m)
	}
	return
}

//
// TagType API.
type TagType struct {
	// hub API client.
	client *Client
}

//
// Create a tag-type.
func (h *TagType) Create(m *api.TagType) (err error) {
	err = h.client.Post(api.TagTypesRoot, m)
	if err == nil {
		Log.Info(
			"Addon created: tag(type).",
			"object",
			m)
	}
	return
}

//
// Get a tag-type by ID.
func (h *TagType) Get(id uint) (m *api.TagType, err error) {
	m = &api.TagType{}
	err = h.client.Get(
		pathlib.Join(
			api.TagTypesRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List tag-types.
func (h *TagType) List() (list []api.TagType, err error) {
	list = []api.TagType{}
	err = h.client.Get(api.TagTypesRoot, &list)
	return
}

//
// Delete a tag-type.
func (h *TagType) Delete(m *api.TagType) (err error) {
	err = h.client.Delete(
		pathlib.Join(
			api.TagTypesRoot,
			strconv.Itoa(int(m.ID))))
	if err == nil {
		Log.Info(
			"Addon deleted: tag(type).",
			"object",
			m)
	}
	return
}
