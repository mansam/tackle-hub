package addon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/konveyor/tackle-hub/api"
	"github.com/konveyor/tackle-hub/task"
	pathlib "path"
	"strconv"
)

//
// Task API.
type Task struct {
	// hub API client.
	client *Client
	// Addon Secret
	secret *task.Secret
	// Task report.
	report api.TaskReport
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
	Log.Info("Addon reported started.")
	return
}

//
// Succeeded report addon succeeded.
func (h *Task) Succeeded() (err error) {
	h.report.Status = task.Succeeded
	err = h.putReport()
	Log.Info("Addon reported: succeeded.")
	return
}

//
// Failed report addon failed.
// The reason can be a printf style format.
func (h *Task) Failed(reason string, x ...interface{}) (err error) {
	h.report.Status = task.Failed
	h.report.Error = fmt.Sprintf(reason, x...)
	err = h.putReport()
	Log.Info(
		"Addon reported: failed.",
		"error",
		h.report.Error)
	return
}

//
// Activity report addon activity.
// The description can be a printf style format.
func (h *Task) Activity(description string, x ...interface{}) (err error) {
	h.report.Activity = fmt.Sprintf(description, x...)
	err = h.putReport()
	Log.Info(
		"Addon reported: activity.",
		"activity",
		h.report.Activity)
	return
}

//
// Total report addon total items.
func (h *Task) Total(n int) (err error) {
	h.report.Total = n
	err = h.putReport()
	Log.Info(
		"Addon updated: total.",
		"total",
		h.report.Total)
	return
}

//
// Increment report addon completed (+1) items.
func (h *Task) Increment() (err error) {
	h.report.Completed++
	err = h.putReport()
	Log.Info(
		"Addon updated: total.",
		"total",
		h.report.Total)
	return
}

//
// Completed report addon completed (N) items.
func (h *Task) Completed(n int) (err error) {
	h.report.Completed = n
	err = h.putReport()
	Log.Info("Addon reported: completed.")
	return
}

//
// postReport creates/updates the task report.
func (h *Task) postReport() (err error) {
	taskID := strconv.Itoa(int(h.secret.Hub.Task))
	err = h.client.Post(
		pathlib.Join(
			api.TasksRoot,
			taskID,
			"report"),
		&h.report)
	if errors.Is(err, &ConflictError{}) {
		err = h.putReport()
	}
	return
}

func (h *Task) putReport() (err error) {
	taskID := strconv.Itoa(int(h.secret.Hub.Task))
	err = h.client.Put(
		pathlib.Join(
			api.TasksRoot,
			taskID,
			"report"),
		&h.report)
	return
}
