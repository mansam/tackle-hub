package model

//
// Proxy configuration.
// kind = (http|https)
type Proxy struct {
	Model
	Kind       string `json:"kind" gorm:"uniqueIndex"`
	Host       string `json:"host" gorm:"not null"`
	Port       int    `json:"port"`
	IdentityID uint   `json:"identity" gorm:"index"`
}
