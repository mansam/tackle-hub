package models

type JobFunctionBinding struct {
	Resource
	Name              string `json:"name"`
	BusinessServiceID string
	BusinessService   BusinessService `json:"business_service"`
	JobFunctionID     string
	JobFunction       JobFunction `json:"job_function"`
	UserID            string
	User              Stakeholder `json:"user"`
	GroupID           string
	Group             Group `json:"group"`
}
