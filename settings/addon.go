package settings

import (
	"net/url"
	"os"
)

const (
	EnvAddonSecretPath  = "ADDON_SECRET_PATH"
	EnvHubBaseURL = "HUB_BASE_URL"
)

//
// Addon settings.
type Addon struct {
	// Hub settings.
	Hub struct {
		// URL for the hub API.
		URL string
	}
	// Shared Secret settings.
	Secret struct {
		// Path to the mounted secret.
		Path string
	}
}

func (r *Addon) Load() (err error) {
	var found bool
	r.Hub.URL, found = os.LookupEnv(EnvHubBaseURL)
	if !found {
		r.Hub.URL = "http://localhost:8080"
	}
	_, err = url.Parse(r.Hub.URL)
	if err != nil {
		panic(err)
	}
	r.Secret.Path, found = os.LookupEnv(EnvAddonSecretPath)
	if !found {
		r.Secret.Path = "/tmp/hub/secret.json"
	}

	return
}
