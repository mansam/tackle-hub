GOBIN ?= ${GOPATH}/bin

# Build ALL commands.
cmd: hub addon

# Run go fmt against code
fmt:
	go fmt ./cmd/... ./api/... ./model/...

# Run go vet against code
vet:
	go vet ./cmd/... ./api/... ./model/...

# Build hub
hub: generate fmt vet
	go build -o bin/hub github.com/konveyor/tackle-hub/cmd

# Build manager binary with compiler optimizations disabled
debug: fmt vet
	go build -o bin/hub -gcflags=all="-N -l" github.com/konveyor/tackle-hub/cmd

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./cmd/main.go

# Generate manifests e.g. CRD, Webhooks
manifests: controller-gen
	${CONTROLLER_GEN} ${CRD_OPTIONS} \
		crd rbac:roleName=manager-role \
		paths="./..." output:crd:artifacts:config=generated/crd/bases output:crd:dir=generated/crd

# Generate code
generate: controller-gen
	${CONTROLLER_GEN} object:headerFile="./generated/boilerplate" paths="./..."

# Find or download controller-gen.
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.5.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# Build SAMPLE ADDON
addon: fmt vet
	go build -o bin/addon github.com/konveyor/tackle-hub/hack/cmd/addon
