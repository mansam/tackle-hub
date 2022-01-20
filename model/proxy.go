package model

import (
	"encoding/json"
	"github.com/konveyor/tackle-hub/encryption"
	"gorm.io/gorm"
)

//
// Proxy configuration.
// kind = (http|https)
type Proxy struct {
	Model
	Kind       string `json:"kind" gorm:"uniqueIndex"`
	Host       string `json:"host" gorm:"not null"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	IdentityID uint   `json:"identity"`
	Encrypted  string `json:"encrypted"`
}

//
// Encrypt sensitive fields.
func (r *Proxy) Encrypt(passphrase string) (err error) {
	aes := encryption.New(passphrase)
	encrypted, err := r.encrypted(passphrase)
	if err != nil {
		return
	}
	if r.User != "" {
		encrypted.User = r.User
		r.User = ""
	}
	if r.Password != "" {
		encrypted.Password = r.Password
		r.Password = ""
	}
	b, err := json.Marshal(encrypted)
	if err != nil {
		return
	}
	r.Encrypted, err = aes.Encrypt(string(b))
	return
}

//
// Decrypt sensitive fields.
func (r *Proxy) Decrypt(passphrase string) (err error) {
	encrypted, err := r.encrypted(passphrase)
	if err != nil {
		return
	}
	r.User = encrypted.User
	r.Password = encrypted.Password

	return
}

//
// BeforeSave ensure encrypted.
func (r *Proxy) BeforeSave(tx *gorm.DB) (err error) {
	err = r.Encrypt(Settings.Encryption.Passphrase)
	return
}

//
// encrypted returns the encrypted identity.
func (r *Proxy) encrypted(passphrase string) (encrypted *Identity, err error) {
	aes := encryption.New(passphrase)
	encrypted = &Identity{}
	if r.Encrypted != "" {
		var dj string
		dj, err = aes.Decrypt(r.Encrypted)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(dj), encrypted)
		if err != nil {
			return
		}
	}
	return
}
