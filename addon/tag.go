package addon

import (
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
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
func (h *Tag) Create(m *model.Tag) (err error) {
	err = h.client.Post(api.TagsRoot, m)
	return
}

//
// Get a tag by ID.
func (h *Tag) Get(id uint) (m *model.Tag, err error) {
	m = &model.Tag{}
	err = h.client.Get(
		pathlib.Join(
			api.TagsRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List tags.
func (h *Tag) List() (list []model.Tag, err error) {
	list = []model.Tag{}
	err = h.client.Get(api.TagsRoot, &list)
	return
}

//
// Delete a tag.
func (h *Tag) Delete(m *model.Tag) (err error) {
	err = h.client.Delete(
		pathlib.Join(
			api.TagsRoot,
			strconv.Itoa(int(m.ID))))
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
func (h *TagType) Create(m *model.TagType) (err error) {
	err = h.client.Post(api.TagTypesRoot, m)
	return
}

//
// Get a tag-type by ID.
func (h *TagType) Get(id uint) (m *model.TagType, err error) {
	m = &model.TagType{}
	err = h.client.Get(
		pathlib.Join(
			api.TagTypesRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List tag-types.
func (h *TagType) List() (list []model.TagType, err error) {
	list = []model.TagType{}
	err = h.client.Get(api.TagTypesRoot, &list)
	return
}

//
// Delete a tag-type.
func (h *TagType) Delete(m *model.TagType) (err error) {
	err = h.client.Delete(
		pathlib.Join(
			api.TagTypesRoot,
			strconv.Itoa(int(m.ID))))
	return
}
