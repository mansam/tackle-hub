package model

import (
	"encoding/json"
	"github.com/konveyor/tackle-hub/encryption"
	"gorm.io/gorm"
)

type Identity struct {
	Model
	Kind         string `json:"kind" gorm:"not null;uniqueIndex:Identity_A"`
	Name         string `json:"name" gorm:"not null"`
	Description  string `json:"description"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Key          string `json:"key"`
	Settings     string `json:"settings"`
	Encrypted    string `json:"encrypted"`
	RepositoryID uint   `json:"repository" gorm:"index;uniqueIndex:Identity_A"`
}

//
// Encrypt sensitive fields.
func (r *Identity) Encrypt(passphrase string) (err error) {
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
	if r.Key != "" {
		encrypted.Key = r.Key
		r.Key = ""
	}
	if r.Settings != "" {
		encrypted.Settings = r.Settings
		r.Settings = ""
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
func (r *Identity) Decrypt(passphrase string) (err error) {
	encrypted, err := r.encrypted(passphrase)
	if err != nil {
		return
	}
	r.User = encrypted.User
	r.Password = encrypted.Password
	r.Key = encrypted.Key
	r.Settings = encrypted.Settings
	return
}

//
// BeforeSave ensure encrypted.
func (r *Identity) BeforeSave(tx *gorm.DB) (err error) {
	err = r.Encrypt(Settings.Encryption.Passphrase)
	return
}

//
// encrypted returns the encrypted identity.
func (r *Identity) encrypted(passphrase string) (encrypted *Identity, err error) {
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
