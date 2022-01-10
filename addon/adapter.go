/*
Tackle hub/addon integration.
 */

package addon

import (
	"encoding/json"
	"github.com/konveyor/tackle-hub/settings"
	"github.com/konveyor/tackle-hub/task"
	"net/http"
	"os"
)

var Settings = settings.Settings

//
// Addon An addon adapter configured for a task execution.
var Addon *Adapter

func init() {
	err := Settings.Load()
	if err != nil {
		panic(err)
	}

	Addon = newAdapter()
}

//
// The Adapter provides hub/addon integration.
type Adapter struct {
	Task
	// Application API.
	Application Application
	// Bucket API.
	Bucket Bucket
	// Identity API.
	Identity Identity
	// TagType API.
	TagType TagType
	// Tag API.
	Tag Tag
	// client A REST client.
	client *Client
}

//
// Client provides the REST client.
func (h *Adapter) Client() *Client {
	return h.client
}

//
// newAdapter builds a new Addon Adapter object.
func newAdapter() (adapter *Adapter) {
	//
	// Load secret.
	secret := &task.Secret{}
	b, err := os.ReadFile(Settings.Addon.Secret.Path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, secret)
	if err != nil {
		panic(err)
	}
	//
	// Build REST client.
	client := &Client{
		baseURL: Settings.Addon.Hub.URL,
		http: &http.Client{},
	}
	//
	// Build Adapter.
	adapter = &Adapter{
		Task: Task{
			client: client,
			secret: secret,
		},
		Application: Application{
			client: client,
		},
		Bucket: Bucket{
			client: client,
		},
		Identity: Identity{
			client: client,
		},
		TagType: TagType{
			client: client,
		},
		Tag: Tag{
			client: client,
		},
		client: client,
	}

	return
}
