package main

import (
	"bytes"
	"fmt"
	hub "github.com/konveyor/tackle-hub/addon"
	"os"
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
	fmt.Printf("Started.")
	_ = addon.Started()
	paths, err := list()
	if err != nil {
		_ = addon.Failed(err.Error())
		os.Exit(-1)
	}
	upload(paths)
	_ = addon.Activity("done")
	_ = addon.Succeeded()
}

//
// upload artifacts.
func upload(paths []string) {
	d := &Data{}
	_ = addon.DataWith(d)
	_ = addon.Total(len(paths))
	for _, p := range paths {
		_ = addon.Activity("uploading: ", p)
		_ = addon.Upload(d.Application, Kind, p)
		pause()
		_ = addon.Completed(1)
	}
}

//
// pause
func pause() {
	time.Sleep(time.Second * 2)
}

//
// list directory.
func list() (paths []string, err error) {
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
