package repository

import (
	"kyberswap_user_monitor/internal/pkg/domain/entity"
)

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination mock_$GOFILE

// TransactionRepository represent repository of the transaction
// Expect implementation by the infrastructure layer
type TransactionRepository interface {
	GetByTxHash(txHash string) (*entity.Transaction, error)
	GetByUser(user string) ([]*entity.Transaction, error)
	Save(user *entity.Transaction) error
	SaveBatch(txBatch []*entity.Transaction) error
}
