package addon

import (
	"errors"
	"github.com/konveyor/tackle-hub/api"
	"net/http"
	pathlib "path"
	"strconv"
)

//
// Repository API.
type Repository struct {
	// hub API client.
	client *Client
}

//
// Get a repository by ID.
func (h *Repository) Get(id uint) (m *api.Repository, err error) {
	m = &api.Repository{}
	err = h.client.Get(
		pathlib.Join(
			api.RepositoriesRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// ByApplication a repository by application ID.
// id - The application ID.
func (h *Repository) ByApplication(id uint) (m *api.Repository, err error) {
	list := []api.Repository{}
	err = h.client.Get(
		pathlib.Join(
			api.AppRepositoriesRoot,
			strconv.Itoa(int(id))),
		&list)
	if err != nil {
		return
	}
	if len(list) < 1 {
		err = errors.New(http.StatusText(http.StatusNotFound))
	}
	m = &list[0]
	return
}

//
// List repositories.
func (h *Repository) List() (list []api.Repository, err error) {
	list = []api.Repository{}
	err = h.client.Get(api.RepositoriesRoot, &list)
	return
}
