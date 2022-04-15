package repository

import "kyberswap_user_monitor/internal/pkg/domain/entity"

type StatRepository interface {
	GetLastStatByChainIDAndType(chainID uint, statType string) (*entity.Stat, error)
	Save(user *entity.Stat) error
}