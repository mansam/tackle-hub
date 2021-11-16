package models

type RoleBinding struct {
	Resource
	Name          string `json:"name"`
	ApplicationID string
	Application   Application `json:"application"`
	RoleID        string
	Role          Role `json:"role"`
	UserID        string
	User          Stakeholder `json:"user"`
	GroupID       string
	Group         Group `json:"group"`
}
