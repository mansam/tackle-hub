/*
TEST addon adapter.
This is an example of an addon adapter that lists files
and creates an application bucket for each. Error handling is
deliberately minimized to reduce code clutter.
 */
package main

import (
	"bytes"
	"errors"
	"fmt"
	hub "github.com/konveyor/tackle-hub/addon"
	"github.com/konveyor/tackle-hub/api"
	"os"
	"os/exec"
	pathlib "path"
	"strconv"
	"strings"
	"time"
)

var (
	// addon adapter.
	addon = hub.Addon
)

//
// main
func main() {
	var err error
	//
	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)
	//
	// Error handler.
	defer func() {
		if err != nil {
			fmt.Printf("Addon failed: %s\n", err.Error())
			_ = addon.Failed(err.Error())
			d.delay()
			os.Exit(1)
		}
	}()
	//
	// Task update: The addon has started.
	// This MUST be called before reporting any
	// other progress.
	_ = addon.Started()
	//
	// Find files.
	paths, _ := find(d.Path, 25)
	//
	// Create bucket.
	err = createBucket(d, paths)
	if err != nil {
		return
	}
	//
	// Task update: The addon has succeeded.
	_ = addon.Succeeded()
}

//
// createBucket builds and populates the bucket.
func createBucket(d *Data, paths []string) (err error) {
	//
	// Task update: Update the task with total number of
	// items to be processed by the addon.
	err = addon.Total(len(paths))
	if err != nil {
		return
	}
	bucket := &api.Bucket{}
	bucket.CreateUser = "addon"
	bucket.Name = "Listing"
	bucket.ApplicationID = d.Application
	err = addon.Bucket.Create(bucket)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = addon.Bucket.Delete(bucket)
		}
	}()
	//
	// Write files.
	for _, p := range paths {
		var b []byte
		//
		// Read file.
		b, err = os.ReadFile(p)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				continue
			}
			return
		}
		//
		// Task update: The current addon activity.
		target := pathlib.Join(
			bucket.Path,
			pathlib.Base(p))
		err = addon.Activity("writing: %s", target)
		if err != nil {
			return
		}
		//
		// Write file.
		err = os.WriteFile(
				target,
				b,
				0644)
		if err != nil {
			return
		}
		time.Sleep(time.Second)
		//
		// Task update: Increment the number of completed
		// items processed by the addon.
		err = addon.Increment()
		if err != nil {
			return
		}
	}
	//
	// Build the index.
	err = buildIndex(bucket)
	if err != nil {
		return
	}
	//
	// Task update: update the current addon activity.
	err = addon.Activity("done")
	return
}

//
// Build index.html
func buildIndex(bucket *api.Bucket) (err error) {
	err = addon.Activity("Building index.")
	time.Sleep(time.Second)
	dir := bucket.Path
	path := pathlib.Join(dir, "index.html")
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	body := []string{"<ul>"}
	list, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, name := range list {
		body = append(
			body,
			"<li><a href=\""+name.Name()+"\">"+name.Name()+"</a>")
	}

	body = append(body, "</ul>")

	_, _ = f.WriteString(strings.Join(body, "\n"))

	return
}

//
// find files.
func find(path string, max int) (paths []string, err error) {
	fmt.Printf("Listing: %s\n", path)
	cmd := exec.Command(
		"find",
		path,
		"-maxdepth",
		"1",
		"-type",
		"f",
		"-readable")
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
	if len(paths) > max {
		paths = paths[:max]
	}

	fmt.Printf("Listed: %v\n", paths)

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
	tag := &api.Tag{}
	tag.Name = "MyTag"
	tag.TagType.ID = 1
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
	// Delay on error (minutes).
	Delay int `json:"delay"`
}

//
// Delay as specified.
func (d *Data) delay() {
	if d.Delay > 0 {
		duration := time.Minute * time.Duration(d.Delay)
		time.Sleep(duration)
	}
}
