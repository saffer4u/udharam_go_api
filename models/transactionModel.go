package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Amount      float32
	Title       string
	LedgerRefer uint
	UserRefer   uint
}
