/*
Tackle hub/addon integration.
 */

package addon

import (
	"encoding/json"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/settings"
	"github.com/konveyor/tackle-hub/task"
	"net/http"
	"os"
	pathlib "path"
	"strconv"
	"strings"
)

var Settings = settings.Settings

//
// Addon An addon adapter configured for a task execution.
var Addon *Adapter

func init() {
	_ = Settings.Load()
	Addon = newAdapter()
}

//
// The Adapter provides hub/addon integration.
type Adapter struct {
	// baseURL for the API.
	baseURL string
	// secret
	secret task.Secret
	// client A REST client.
	client *Client
	// Task report.
	report model.TaskReport
}

//
// Client provides the REST client.
func (h *Adapter) Client() *Client {
	return h.client
}

//
// Data returns the addon data.
func (h *Adapter) Data() (d map[string]interface{}) {
	d = h.secret.Addon.(map[string]interface{})
	return
}

//
// Application provides the application
// to be processed by the addon.
func (h *Adapter) Application(id int) (m *model.Application, err error) {
	m = &model.Application{}
	err = h.client.Get(
		pathlib.Join(
			api.ApplicationsRoot,
			strconv.Itoa(id)),
		m)
	return
}

//
// Started report addon started.
func (h *Adapter) Started() (err error) {
	h.report.Status = task.Running
	err = h.postReport()
	return
}

//
// Succeeded report addon succeeded.
func (h *Adapter) Succeeded() (err error) {
	h.report.Status = task.Succeeded
	err = h.putReport()
	return
}

//
// Failed report addon failed.
func (h *Adapter) Failed(reason string) (err error) {
	h.report.Status = task.Failed
	h.report.Error = reason
	err = h.putReport()
	return
}

//
// Activity report addon activity.
func (h *Adapter) Activity(word ...string) (err error) {
	h.report.Activity = strings.Join(word, " ")
	err = h.putReport()
	return
}

//
// Total report addon total items.
func (h *Adapter) Total(n int) (err error) {
	h.report.Total = n
	err = h.postReport()
	return
}

//
// Completed report addon completed (N) items.
func (h *Adapter) Completed(n int) (err error) {
	h.report.Completed += n
	err = h.putReport()
	return
}

//
// Upload an artifact.
func (h *Adapter) Upload(application uint, kind string, path string) (err error) {
	artifact := &model.Artifact{}
	artifact.CreateUser = "addon"
	artifact.Name = pathlib.Base(path)
	artifact.ApplicationID = application
	artifact.Kind = kind
	err = h.client.Post(api.ArtifactsRoot, artifact)
	return
}

func (h *Adapter) postReport() (err error) {
	taskID := strconv.Itoa(int(h.secret.Hub.Task))
	h.report.CreateUser = "addon"
	err = h.client.Post(
		pathlib.Join(
			api.TasksRoot,
			taskID,
			"report"),
		&h.report)
	return
}

func (h *Adapter) putReport() (err error) {
	taskID := strconv.Itoa(int(h.secret.Hub.Task))
	h.report.UpdateUser = "addon"
	err = h.client.Put(
		pathlib.Join(
			api.TasksRoot,
			taskID,
			"report"),
		&h.report)
	return
}

//
// newHub builds a new Addon Adapter object.
func newAdapter() (adapter *Adapter) {
	adapter = &Adapter{}
	// base URL
	adapter.baseURL = Settings.Addon.Hub.URL
	// Load secret.
	b, err := os.ReadFile(Settings.Addon.Secret.Path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &adapter.secret)
	if err != nil {
		panic(err)
	}
	// REST client.
	adapter.client = &Client{
		baseURL: adapter.baseURL,
		http: &http.Client{},
	}

	return
}
