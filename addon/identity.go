package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
	"strconv"
)

//
// Identity API.
type Identity struct {
	// hub API client.
	client *Client
}

//
// Get an identity by ID.
func (h *Identity) Get(id uint) (m *api.Identity, err error) {
	m = &api.Identity{}
	err = h.client.Get(
		pathlib.Join(
			api.IdentitiesRoot,
			strconv.Itoa(int(id))),
		m)
	if err != nil {
		return
	}
	err = m.Decrypt(Addon.secret.Hub.Encryption.Passphrase)
	return
}

//
// List identities.
func (h *Identity) List() (list []api.Identity, err error) {
	list = []api.Identity{}
	err = h.client.Get(api.IdentitiesRoot, &list)
	if err != nil {
		return
	}
	for i := range list {
		m := &list[i]
		err = m.Decrypt(Addon.secret.Hub.Encryption.Passphrase)
		if err != nil {
			return
		}
	}
	return
}
