package addon

import (
	"encoding/json"
	"fmt"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/task"
	pathlib "path"
	"strconv"
	"strings"
)

//
// Task API.
type Task struct {
	// hub API client.
	client *Client
	// Addon Secret
	secret *task.Secret
	// Task report.
	report model.TaskReport

}

//
// Data returns the addon data.
func (h *Task) Data() (d map[string]interface{}) {
	d = h.secret.Addon.(map[string]interface{})
	return
}

//
// DataWith populates the addon data object.
func (h *Task) DataWith(object interface{}) (err error) {
	b, _ := json.Marshal(h.secret.Addon)
	err = json.Unmarshal(b, object)
	return
}

//
// Started report addon started.
func (h *Task) Started() (err error) {
	h.report.Status = task.Running
	err = h.postReport()
	return
}

//
// Succeeded report addon succeeded.
func (h *Task) Succeeded() (err error) {
	h.report.Status = task.Succeeded
	err = h.putReport()
	return
}

//
// Failed report addon failed.
// The reason can be a printf style format.
func (h *Task) Failed(reason string, x ...interface{}) (err error) {
	h.report.Status = task.Failed
	h.report.Error = fmt.Sprintf(reason, x...)
	err = h.putReport()
	return
}

//
// Activity report addon activity.
func (h *Task) Activity(word ...string) (err error) {
	h.report.Activity = strings.Join(word, " ")
	err = h.putReport()
	return
}

//
// Total report addon total items.
func (h *Task) Total(n int) (err error) {
	h.report.Total = n
	err = h.postReport()
	return
}

//
// Increment report addon completed (+1) items.
func (h *Task) Increment() (err error) {
	h.report.Completed++
	err = h.putReport()
	return
}

//
// Completed report addon completed (N) items.
func (h *Task) Completed(n int) (err error) {
	h.report.Completed = n
	err = h.putReport()
	return
}

func (h *Task) postReport() (err error) {
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

func (h *Task) putReport() (err error) {
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
