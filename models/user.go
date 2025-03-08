package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	Password   string
	MFASecret  string `gorm:"column:mfa_secret"`
	MFAEnabled bool   `gorm:"column:mfa_enabled"`
}
