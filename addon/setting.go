package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
)

//
// Setting API.
type Setting struct {
	// hub API client.
	client *Client
}

//
// Get a setting by key.
func (h *Setting) Get(key string) (v interface{}, err error) {
	r := &api.Setting{}
	err = h.client.Get(
		pathlib.Join(
			api.SettingsRoot,
			key),
		r)
	v = r.Value
	return
}
