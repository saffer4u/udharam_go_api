package models

import "gorm.io/gorm"

type Ledger struct {
	gorm.Model
	LedgerName  string
	UserRefer   uint
	Transaction []Transaction `gorm:"foreignKey:LedgerRefer"`
}
