package settings

import (
	"net/url"
	"os"
)

const (
	AddonSecret  = "SECRET"
	EnvBaseURL = "HUB"
)

type Addon struct {
	API struct {
		URL string
	}
	Secret struct {
		Path string
	}
}

func (r *Addon) Load() (err error) {
	var found bool
	r.API.URL, found = os.LookupEnv(EnvBaseURL)
	if !found {
		r.API.URL = "http://localhost:8080"
	}
	_, err = url.Parse(r.API.URL)
	if err != nil {
		panic(err)
	}
	r.Secret.Path, found = os.LookupEnv(AddonSecret)
	if !found {
		r.Secret.Path = "/tmp/hub/secret.json"
	}

	return
}
