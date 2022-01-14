package model

const (
	HubSettingsKey = "hub"
	MvnSettingsKey = "mvn"
)

//
// Setting stored in the DB.
type Setting struct {
	Model
	Key string `json:"key"`
	Value string `json:"value"`
}

//
// HubSettings hib settings.
type HubSettings struct {
	DB struct {
		Seeded bool `json:"seeded"`
	} `json:"db"`
}

type MvnSettings struct {
	AllowInsecure bool `json:"allowInsecure"`
	UpdateForced bool `json:"updateForced"`
}