package model

import (
	"encoding/json"
	"gorm.io/gorm"
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
	Data       interface{} `json:"data" gorm:"-"`
	Data_      string      `json:"-"`
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

func (m *Task) BeforeSave(db *gorm.DB) (err error) {
	if m.Data == nil {
		m.Data = struct{}{}
	}
	b, err := json.Marshal(m.Data)
	m.Data_ = string(b)
	return
}

func (m *Task) AfterFind(db *gorm.DB) (err error) {
	if len(m.Data_) == 0 {
		m.Data_ = "{}"
	}
	b := []byte(m.Data_)
	err = json.Unmarshal(b, &m.Data)
	return
}
