package addon

import (
	"github.com/konveyor/tackle-hub/api"
	pathlib "path"
	"strconv"
)

//
// Bucket API.
type Bucket struct {
	// hub API client.
	client *Client
}

//
// Create a bucket.
func (h *Bucket) Create(m *api.Bucket) (err error) {
	err = h.client.Post(api.BucketsRoot, m)
	if err == nil {
		Log.Info(
			"Addon created: bucket.",
			"object",
			m)
	}
	return
}

//
// Get a bucket by ID.
func (h *Bucket) Get(id uint) (m *api.Bucket, err error) {
	m = &api.Bucket{}
	err = h.client.Get(
		pathlib.Join(
			api.BucketsRoot,
			strconv.Itoa(int(id))),
		m)
	return
}

//
// List buckets.
func (h *Bucket) List() (list []api.Bucket, err error) {
	list = []api.Bucket{}
	err = h.client.Get(api.BucketsRoot, &list)
	return
}

//
// Delete an bucket.
func (h *Bucket) Delete(m *api.Bucket) (err error) {
	err = h.client.Delete(
		pathlib.Join(
			api.BucketsRoot,
			strconv.Itoa(int(m.ID))))
	if err == nil {
		Log.Info(
			"Addon deleted: bucket.",
			"object",
			m)
	}
	return
}
