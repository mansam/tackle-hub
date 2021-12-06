package settings

import "os"

const (
	EnvDbPath  = "DB_PATH"
)

type Hub struct {
	DB struct {
		Path string
	}
}

func (r *Hub) Load() (err error) {
	var found bool
	r.DB.Path, found = os.LookupEnv(EnvDbPath)
	if !found {
		r.DB.Path = "/tmp/tackle.db"
	}

	return
}