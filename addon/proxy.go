package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
	"strconv"
)

//
// Proxy API.
type Proxy struct {
	// hub API client.
	client *Client
}

//
// Get a proxy by ID.
func (h *Proxy) Get(id uint) (m *api.Proxy, err error) {
	m = &api.Proxy{}
	err = h.client.Get(
		pathlib.Join(
			api.ProxiesRoot,
			strconv.Itoa(int(id))),
		m)
	if err != nil {
		return
	}
	err = m.Decrypt(Addon.secret.Hub.Encryption.Passphrase)
	return
}

//
// List proxies.
func (h *Proxy) List() (list []api.Proxy, err error) {
	list = []api.Proxy{}
	err = h.client.Get(api.ProxiesRoot, &list)
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

//
// Update a proxy by ID.
func (h *Proxy) Update(m *api.Proxy) (err error) {
	err = h.client.Put(
		pathlib.Join(
			api.ProxiesRoot,
			strconv.Itoa(int(m.ID))),
		m)
	return
}
