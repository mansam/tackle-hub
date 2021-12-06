package settings

import "os"

const (
	EnvNamespace = "NAMESPACE"
	EnvDbPath  = "DB_PATH"
)

type Hub struct {
	// k8s namespace.
	Namespace string
	// DB settings.
	DB struct {
		Path string
	}
}

func (r *Hub) Load() (err error) {
	var found bool
	r.DB.Path, found = os.LookupEnv(EnvNamespace)
	if !found {
		r.DB.Path = "tackle-hub"
	}
	r.DB.Path, found = os.LookupEnv(EnvDbPath)
	if !found {
		r.DB.Path = "/tmp/tackle.db"
	}

	return
}