package model

import (
	"time"
)

type TaskReport struct {
	Model
	Status    string
	Error     string
	Total     int
	Completed int
	Activity  string
	TaskID    uint
	Task      *Task
}

type Task struct {
	Model
	Name       string
	Image      string
	Addon      string
	Isolated   bool
	Data       JSON
	Started    *time.Time
	Terminated *time.Time
	Status     string
	Error      string
	Job        string
	Report     *TaskReport `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Task) Reset() {
	m.Started = nil
	m.Terminated = nil
	m.Report = nil
	m.Status = ""
}
