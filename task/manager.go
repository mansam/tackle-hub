package task

import (
	"context"
	"encoding/json"
	crd "github.com/konveyor/tackle-hub/k8s/api/tackle/v1alpha1"
	"github.com/konveyor/tackle-hub/model"
	"github.com/konveyor/tackle-hub/settings"
	"gorm.io/gorm"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"time"
)

const (
	Succeeded = "Succeeded"
	Failed = "Failed"
	Running = "Running"
)

const (
	AppLabel = "tackle-hub"
)

var Settings = &settings.Settings

//
// Manager provides task management.
type Manager struct {
	// DB
	DB *gorm.DB
	// k8s client.
	Client client.Client
}

//
// Run the manager.
func (m *Manager) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second)
				_ = m.runPending()
				_ = m.reflectRunning()
			}
		}
	}()
}

//
// runPending starts (run) pending tasks.
func (m *Manager) runPending() (err error) {
	list := []model.Task{}
	result := m.DB.Find(&list, "status", "")
	if result.Error != nil {
		err = result.Error
		return
	}
	for _, pending := range list {
		task := Task{
			Task: &pending,
			client: m.Client,
		}
		_ = task.Run()
		_ = m.DB.Save(&pending)
	}

	return
}

//
// reflectRunning updates running tasks to reflect job status.
func (m *Manager) reflectRunning() (err error) {
	list := []model.Task{}
	result := m.DB.Find(&list, "status", Running)
	if result.Error != nil {
		err = result.Error
		return
	}
	for _, running := range list {
		task := Task{
			Task: &running,
			client: m.Client,
		}
		err := task.Reflect()
		if err != nil {
			continue
		}
		_ = m.DB.Save(&running)
	}

	return
}

//
// Task is an runtime task.
type Task struct {
	// model.
	*model.Task
	// k8s client.
	client client.Client
}

//
// Run the specified task.
func (r *Task) Run() (err error) {
	defer func() {
		if err != nil {
			r.Error = err.Error()
			r.Status = Failed
		}
	}()
	addon, err := r.findAddon(r.Addon)
	if err != nil {
		return
	}
	r.Image = addon.Spec.Image
	secret := r.secret()
	err = r.client.Create(context.TODO(), &secret)
	if err != nil {
		return
	}
	job := r.job(&secret)
	err = r.client.Create(context.TODO(), &job)
	if err != nil {
		return
	}
	mark := time.Now()
	r.Started = &mark
	r.Status = Running
	r.Job = path.Join(
		job.Namespace,
		job.Name)
	return
}

//
// Reflect finds the associated job and updates the task status.
func (r *Task) Reflect() (err error) {
	job := &batch.Job{}
	err = r.client.Get(
		context.TODO(),
		client.ObjectKey{
			Namespace: path.Dir(r.Job),
			Name: path.Base(r.Job),
		},
		job)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Run()
		}
		return
	}
	mark := time.Now()
	status := job.Status
	for _, cnd := range status.Conditions {
		if cnd.Type == batch.JobFailed {
			r.Status = Failed
			r.Terminated = &mark
			r.Error = "job failed."
			return
		}
		if status.Succeeded > 0 {
			r.Status = Succeeded
			r.Terminated = &mark
		}
	}

	return
}

//
// findAddon by name.
func (r *Task) findAddon(name string) (addon *crd.Addon, err error) {
	addon = &crd.Addon{}
	err = r.client.Get(
		context.TODO(),
		client.ObjectKey{
			Namespace: Settings.Hub.Namespace,
			Name: name,
		},
		addon)
	if err != nil {
		return
	}

	return
}

//
// job build the Job.
func (r *Task) job(secret *core.Secret) (job batch.Job) {
	template := r.template(secret)
	job = batch.Job{
		Spec: batch.JobSpec{Template: template},
		ObjectMeta: meta.ObjectMeta{
			Namespace: Settings.Hub.Namespace,
			GenerateName: strings.ToLower(r.Name) + "-",
			Labels: r.labels(),
		},
	}

	return
}

//
// template builds a Job template.
func (r *Task) template(secret *core.Secret) (template core.PodTemplateSpec) {
	template = core.PodTemplateSpec{
		Spec: core.PodSpec{
			Affinity: r.affinity(),
			RestartPolicy: "OnFailure",
			Containers: []core.Container{
				r.container(),
			},
			Volumes: []core.Volume{
				{
					Name: "secret",
					VolumeSource: core.VolumeSource{
						Secret: &core.SecretVolumeSource{
							SecretName: secret.Name,
						},
					},
				},
				{
					Name: "bucket",
					VolumeSource: core.VolumeSource{
						PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
							ClaimName: Settings.Hub.Bucket.PVC,
						},
					},
				},
			},
		},
	}

	return
}

//
// container builds the job container.
func (r *Task) container() (container core.Container) {
	container = core.Container{
		Name:  "main",
		Image: r.Image,
		Env: []core.EnvVar {
			{
				Name: settings.EnvHubBaseURL,
				Value: Settings.Addon.Hub.URL,
			},
			{
				Name: settings.EnvAddonSecretPath,
				Value: Settings.Addon.Secret.Path,
			},
			{
				Name: settings.EnvBucketPath,
				Value: Settings.Hub.Bucket.Path,
			},
		},
		VolumeMounts: []core.VolumeMount{
			{
				Name:      "secret",
				MountPath: path.Dir(Settings.Addon.Secret.Path),
			},
			{
				Name:      "bucket",
				MountPath: Settings.Hub.Bucket.Path,
			},
		},
	}

	return
}

//
// Affinity.
func (r *Task) affinity() (affinity *core.Affinity) {
	affinity = &core.Affinity{
		PodAffinity: &core.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{
				{
					TopologyKey: "kubernetes.io/hostname",
					LabelSelector: &meta.LabelSelector{
						MatchLabels: map[string]string{
							"app": AppLabel,
						},
					},
				},
			},
		},
	}

	return
}

//
// secret builds the job secret.
func (r *Task) secret() (secret core.Secret) {
	data := Secret{}
	data.Hub.Task = r.Task.ID
	data.Addon = r.Task.Data
	encoded, _ := json.Marshal(data)
	secret = core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: Settings.Hub.Namespace,
			GenerateName: strings.ToLower(r.Name) + "-",
			Labels: r.labels(),
		},
		Data: map[string][]byte{
			path.Base(Settings.Addon.Secret.Path): encoded,
		},
	}

	return
}

//
// labels builds k8s labels.
func (r *Task) labels() map[string]string {
	return map[string]string {
		"Task": strconv.Itoa(int(r.ID)),
	}
}

//
// Secret payload.
type Secret struct {
	Hub struct {
		Token string
		Task uint
	}
	Addon interface{}
}
