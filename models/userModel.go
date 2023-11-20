package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
	Ledger    []Ledger `gorm:"foreignKey:UserRefer"`
}
