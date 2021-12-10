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
	"github.com/konveyor/tackle-hub/model"
	"os/exec"
	"strconv"
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
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("Addon failed: %s", err.Error())
			_ = addon.Failed(err.Error())
		}
	}()
	//
	// Task update: The addon has started.
	// This MUST be called before reporting any
	// other progress.
	_ = addon.Started()
	//
	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)
	//
	// Find files.
	paths, _ := find(d.Path)
	//
	// Upload files and create artifacts.
	err = upload(d, paths)
	if err != nil {
		return
	}
	//
	// Task update: The addon has succeeded.
	_ = addon.Succeeded()
}

//
// upload artifacts.
func upload(d *Data, paths []string) (err error) {
	//
	// Task update: Update the task with total number of
	// items to be processed by the addon.
	err = addon.Total(len(paths))
	if err != nil {
		return
	}
	//
	// Upload artifacts.
	for _, p := range paths {
		//
		// Task update: The current addon activity.
		err = addon.Activity("uploading: ", p)
		if err != nil {
			return
		}
		//
		// Upload the file and create an artifact to be
		// associated with the application.
		err = addon.Artifact.Upload(d.Application, Kind, p)
		if err != nil {
			return
		}
		pause()
		//
		// Task update: Increment the number of completed
		// items processed by the addon.
		err = addon.Increment()
		if err != nil {
			return
		}
	}
	//
	// Task update: update the current addon activity.
	err = addon.Activity("done")
	return
}

//
// pause
func pause() {
	time.Sleep(time.Second * 2)
}

//
// find files.
func find(path string) (paths []string, err error) {
	fmt.Printf("Listing: %s", path)
	cmd := exec.Command(
		"find",
		path,
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

	max := 10
	paths = strings.Fields(stdout.String())
	if len(paths) > max {
		paths = paths[:max]
	}

	fmt.Printf("Listed: %v", paths)

	return
}

//
// Tag application.
func tag(d *Data) (err error) {
	//
	// Fetch application.
	application, _ := addon.Application.Get(d.Application)
	//
	// Create tag.
	tag := &model.Tag{
		Name: "MyTag",
		TagTypeID: 1,
	}
	err = addon.Tag.Create(tag)
	if err != nil {
		return
	}
	//
	// append tag.
	application.Tags = append(
		application.Tags,
		strconv.Itoa(int(tag.ID)))
	//
	// Update application.
	err = addon.Application.Update(application)
	return
}

//
// Data Addon data passed in the secret.
type Data struct {
	// Application ID.
	Application uint `json:"application"`
	// Path to be listed.
	Path string `json:"path"`
}
