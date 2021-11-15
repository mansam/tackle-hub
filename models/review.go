package models

type Review struct {
	Resource
	BusinessCriticality uint   `json:"business_criticality" gorm:"notnull" binding:"required"`
	EffortEstimate      string `json:"effort_estimate" gorm:"notnull" binding:"required"`
	ProposedAction      string `json:"proposed_action" gorm:"notnull" binding:"required"`
	WorkPriority        uint   `json:"work_priority" gorm:"notnull" binding:"required"`
	ApplicationID       string `json:"application_id"`
	Application         Application
}
