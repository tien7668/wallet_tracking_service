package entity

import (
	"gorm.io/gorm"
)

const (
	TypeSwap = "swap"
	TypeAdd	= "add"		
	TypeRemove = "remove"
)

type Stat struct {
	ChainID   	uint   `gorm:"primaryKey"`
	CrawlType   string `gorm:"primaryKey;type:varchar(16)"`
	Timestamp 	uint64
}

func StatTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table("stats")
	}
}


func (Stat) TableName() string {
	return "stats"
}