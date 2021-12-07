/*
TEST addon adapter.
This is an example of an addon adapter that lists files in /etc
and creates an application artifact for each. Error handling is
deliberately minimized to reduce code clutter.
 */
package main

import (
	"bytes"
	"fmt"
	hub "github.com/konveyor/tackle-hub/addon"
	"os/exec"
	"strings"
	"time"
)

const (
	Kind = "File"
)

var (
	// addon adapter.
	addon = hub.Addon
)

//
// main
func main() {
	//
	// Task update: The addon has started.
	// This MUST be called before reporting any
	// other progress.
	_ = addon.Started()
	//
	// Find files.
	paths, _ := find()
	//
	// Upload files and create artifacts.
	upload(paths)
	//
	// Task update: The addon has succeeded.
	_ = addon.Succeeded()
}

//
// upload artifacts.
func upload(paths []string) {
	//
	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)
	//
	// Task update: Update the task with total number of
	// items to be processed by the addon.
	_ = addon.Total(len(paths))
	//
	// Upload artifacts.
	for _, p := range paths {
		//
		// Task update: The current addon activity.
		_ = addon.Activity("uploading: ", p)
		//
		// Upload the file and create an artifact to be
		// associated with the application.
		_ = addon.Upload(d.Application, Kind, p)
		pause()
		//
		// Task update: Increment the number of completed
		// items processed by the addon.
		_ = addon.Completed(1)
	}
	//
	// Task update: update the current addon activity.
	_ = addon.Activity("done")
}

//
// pause
func pause() {
	time.Sleep(time.Second * 2)
}

//
// find files.
func find() (paths []string, err error) {
	cmd := exec.Command(
		"find",
		"/etc",
		"-maxdepth",
		"1",
		"-type",
		"f")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}

	paths = strings.Fields(stdout.String())

	fmt.Printf("Listed: %v", paths)

	return
}

//
// Data Addon data passed in the secret.
type Data struct {
	Application uint `json:"application"`
}
