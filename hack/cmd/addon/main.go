package main

import (
	"bytes"
	"fmt"
	hub "github.com/konveyor/tackle-hub/addon"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	Kind = "File"
)

var (
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
// application returns the application ID
// specified in the secret.
func application() (id uint) {
	value := addon.Data()["application"]
	if f, cast := value.(float64); cast {
		id = uint(f)
	}

	return
}

//
// upload artifacts.
func upload(paths []string) {
	_ = addon.Total(len(paths))
	id := application()
	for _, p := range paths {
		_ = addon.Activity("uploading: ", p)
		_ = addon.Upload(id, Kind, p)
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
	dir := "/var/log"
	cmd := exec.Command("ls", dir)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	files := strings.Fields(stdout.String())
	for _, name := range files {
		paths = append(
			paths,
			path.Join(dir, name))
	}

	fmt.Printf("Listed: %v", paths)

	return
}
