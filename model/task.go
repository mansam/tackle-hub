package model

import (
	"time"
)

type TaskReport struct {
	Model
	Status    string `json:"status"`
	Error     string `json:"error"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
	Activity  string `json:"activity"`
	TaskID    uint   `json:"task"`
	Task      *Task  `json:"-"`
}

type Task struct {
	Model
	Name       string      `json:"name"`
	Image      string      `json:"image"`
	Addon      string      `json:"addon"`
	Data       JSON        `json:"data"`
	Started    *time.Time  `json:"started"`
	Terminated *time.Time  `json:"terminated"`
	Status     string      `json:"status"`
	Error      string      `json:"error"`
	Job        string      `json:"job"`
	Report     *TaskReport `json:"report" gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Task) Reset() {
	m.Started = nil
	m.Terminated = nil
	m.Report = nil
	m.Status = ""
}
