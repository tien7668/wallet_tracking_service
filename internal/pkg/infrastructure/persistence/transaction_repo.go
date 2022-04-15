package persistence

import (
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/domain/repository"
	"kyberswap_user_monitor/internal/pkg/state"
)

type TransactionRepoImpl struct {}
var repoImpl *TransactionRepoImpl
func GetTransactionRepoImpl() repository.TransactionRepository {
	if repoImpl != nil {
		return repoImpl
	} 
	return &TransactionRepoImpl{}
}

func (r *TransactionRepoImpl) GetByTxHash(txHash string) (*entity.Transaction, error) {
	db,_ := state.GetDB()
	var tx entity.Transaction
	if err := db.Where(&entity.Transaction{Tx: txHash}).First(&tx).Error; err != nil {
		return nil, err
	} 
	return &tx, nil
} 

func (r *TransactionRepoImpl) GetByUser(user string) ([]*entity.Transaction, error) {
	db,_ := state.GetDB()
	var txs []*entity.Transaction
	if err := db.Where(&entity.Transaction{UserAddress: user}).Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *TransactionRepoImpl) Save(tx *entity.Transaction) error {
	db,_ := state.GetDB()
	if err := db.Save(tx); err != nil {
		return err.Error
	}
	return nil
} 

func (r *TransactionRepoImpl) SaveBatch(list []*entity.Transaction) error {
	db,_ := state.GetDB()
	if len(list) == 0 {
		return nil
	}
	if err := db.Create(list).Error; err != nil {
		return err
	}
	return nil
}