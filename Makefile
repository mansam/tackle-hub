
# Run go fmt against code
fmt:
	go fmt ./db/... ./cmd/... ./handlers/... ./models/...

# Run go vet against code
vet:
	go vet ./db/... ./cmd/... ./handlers/... ./models/...

# Build server binary
server: fmt vet
	go build -o bin/serve github.com/konveyor/tackle-hub/cmd

# Build manager binary with compiler optimizations disabled
debug: fmt vet
	go build -o bin/serve -gcflags=all="-N -l" github.com/konveyor/tackle-hub/cmd

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./cmd/serve.go
