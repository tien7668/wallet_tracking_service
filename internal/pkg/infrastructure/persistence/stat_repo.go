package persistence

import (
	"errors"
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/domain/repository"
	"kyberswap_user_monitor/internal/pkg/state"

	"gorm.io/gorm"
)

type StatRepoImpl struct {}
var statImpl *StatRepoImpl
func GetStatRepoImpl() repository.StatRepository {
	if repoImpl != nil {
		return statImpl
	} 
	return &StatRepoImpl{}
}

func (r *StatRepoImpl) GetLastStatByChainIDAndType(chainID uint, statType string) (*entity.Stat, error) {
	db,_ := state.GetDB()
	ctx := state.GetContext()
	var lastItem *entity.Stat
	query := db.Where(&entity.Stat{ChainID: chainID, CrawlType: statType}).Order("timestamp desc")
	err := query.Take(&lastItem).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound){
			ctx.Errorf("failed to get last stat of chainID: %v and type: %v , err: %v", chainID, statType, err)
			return nil, err
		}else {
			lastItem = &entity.Stat{ChainID: chainID, CrawlType: statType}
			db.Create(lastItem)
			// ctx.Infof("------------------------------ %v %v %v ", chainID, statType, lastItem)
			return lastItem, nil
		}
	} 
	return lastItem, nil

}

func (r *StatRepoImpl) Save(stat *entity.Stat) error {
	db,_ := state.GetDB()
	// ctx := state.GetContext()
	db.Save(stat)
	return nil
} 