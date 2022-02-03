package addon

import (
	"errors"
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

//
// Bool setting value.
func (h *Setting) Bool(key string) (b bool, err error) {
	v, err := h.Get(key)
	if err != nil {
		return
	}
	b, cast := v.(bool)
	if !cast {
		err = errors.New(key + " not <boolean>")
	}
	return
}

//
// Str setting value.
func (h *Setting) Str(key string) (s string, err error) {
	v, err := h.Get(key)
	if err != nil {
		return
	}
	s, cast := v.(string)
	if !cast {
		err = errors.New(key + " not <string>")
	}
	return
}

//
// Int setting value.
func (h *Setting) Int(key string) (n int, err error) {
	v, err := h.Get(key)
	if err != nil {
		return
	}
	n, cast := v.(int)
	if !cast {
		err = errors.New(key + " not <int>")
	}
	return
}
