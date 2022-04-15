package entity

import (
	"gorm.io/gorm"
	"encoding/json"
)

const (
	TxTypeSwap = "swap"
	TxTypeAdd = "add"
	TxTypeRemove = "remove"
)
type Transaction struct {
	Id          uint `gorm:"primaryKey;auto_increment;not_null"`
	ChainID     uint
	Timestamp   uint64
	Pair        string
	Extra		[]byte
	UserAddress string `gorm:"index"`
	BlockNumber uint64
	Tx          string
	LastSync    uint64
	TxType 		string
}

type TransactionJSON struct {
	Id          uint
	ChainID     uint
	Timestamp   uint64
	Pair        string
	UserAddress string
	BlockNumber uint64
	Tx          string
	LastSync    uint64
	TxType 		string
	Extra 		map[string]interface{}
} 

func (t *Transaction) ToJSON() (TransactionJSON){
	m := map[string]interface{}{}
	json.Unmarshal([]byte(t.Extra), &m)
	return TransactionJSON{
		Id 			:t.Id,
		ChainID		:t.ChainID,
		Timestamp	:t.Timestamp,
		Pair		:t.Pair,
		Extra		:m,
		UserAddress	:t.UserAddress,
		BlockNumber	:t.BlockNumber,
		Tx			:t.Tx,
		LastSync	:t.LastSync,
		TxType		:t.TxType,
	}
}

func TransactionTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table("transactions")
	}
}

func (Transaction) TableName() string {
	return "transactions"
}
