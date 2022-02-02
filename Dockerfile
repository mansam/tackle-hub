FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12 as builder
ENV GOPATH=$APP_ROOT
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o hub github.com/konveyor/tackle-hub/cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4
COPY --from=builder /opt/app-root/src/hub /usr/local/bin/tackle-hub
ENTRYPOINT ["/usr/local/bin/tackle-hub"]

LABEL name="konveyor/tackle-hub" \
      description="Konveyor Tackle - Hub" \
      help="For more information visit https://konveyor.io" \
      license="Apache License 2.0" \
      maintainers="jortel@redhat.com,slucidi@redhat.com" \
      summary="Konveyor Tackle - Hub" \
      url="https://quay.io/repository/konveyor/tackle-hub" \
      usage="podman run konveyor/tackle-hub:latest" \
      com.redhat.component="konveyor-tackle-hub-container" \
      io.k8s.display-name="Tackle Hub" \
      io.k8s.description="Konveyor Tackle - Hub" \
      io.openshift.expose-services="" \
      io.openshift.tags="konveyor,tackle,hub" \
      io.openshift.min-cpu="100m" \
      io.openshift.min-memory="350Mi"
