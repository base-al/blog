package models

import (
	"gorm.io/gorm"
)

type Administrator struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Role     string `gorm:"not null"`
}
type Setting struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	Value  string `gorm:"not null"`
	Module string `gorm:"not null"`
}
