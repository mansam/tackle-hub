package addon

import (
	"github.com/konveyor/tackle-hub/api"
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
// List repositories.
func (h *Repository) List() (list []api.Repository, err error) {
	list = []api.Repository{}
	err = h.client.Get(api.RepositoriesRoot, &list)
	return
}
