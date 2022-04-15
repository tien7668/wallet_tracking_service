package application

import (
    "kyberswap_user_monitor/internal/pkg/domain/entity"
    "kyberswap_user_monitor/internal/pkg/domain/repository"
    "kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl"
	"kyberswap_user_monitor/internal/pkg/domain/service/token_info"
    "kyberswap_user_monitor/internal/pkg/state"
    "time"
    // "sync"
)

type TransactionUsecase struct {
    repository.TransactionRepository
    repository.StatRepository
    SwapCrawl tx_crawl.CrawlInterface
    AddLiquidCrawl tx_crawl.CrawlInterface
	RemoveLiquidCrawl tx_crawl.CrawlInterface
	Tokens map[string]token_info.TokenInfo
} 

var transactionUsecase TransactionUsecase

func GetTransactionUsecase() TransactionUsecase{
    return transactionUsecase
}

func InitTransactionUsecase(t TransactionUsecase) {
    transactionUsecase = t
} 


const (
    statsJobIntervalSec = 60
    retryTimeout        = 15
)

func (i TransactionUsecase) Crawl(getLastStat func()(*entity.Stat, error), crawl func(lastStat *entity.Stat)([]*entity.Transaction, uint64, error)) {

	for {
		lastStat, err := getLastStat()
		if err == nil {
			list, lastTimestamp, err := crawl(lastStat)
			if err != nil {
				return 
			} 
			if lastTimestamp > 0 && len(list) > 0 {
				lastStat.Timestamp = lastTimestamp + 1
				err = i.StatRepository.Save(lastStat)
				i.TransactionRepository.SaveBatch(list)	
			}
		}
		time.Sleep(statsJobIntervalSec * time.Second / 10)
	}
}

func (i TransactionUsecase) CrawAll() error {
    ctx := state.GetContext()
    cfg,_ := state.GetConfig()
	ctx.Infof("-------START CRAWL")
	for _, chain := range cfg.Blockchain.Chains {
		chainID := chain.ChainID
		infoApi := chain.API.Price
		go i.Crawl(
			func()(*entity.Stat, error){
				return i.StatRepository.GetLastStatByChainIDAndType(chainID, entity.TypeSwap)
			}, 
			func(lastStat *entity.Stat)([]*entity.Transaction, uint64, error){
				return i.SwapCrawl.Crawl(chainID, infoApi, lastStat.Timestamp)
			},
		)
		go i.Crawl(
			func()(*entity.Stat, error){
				return i.StatRepository.GetLastStatByChainIDAndType(chainID, entity.TypeAdd)
			}, 
			func(lastStat *entity.Stat)([]*entity.Transaction, uint64, error){
				return i.AddLiquidCrawl.Crawl(chainID, infoApi, lastStat.Timestamp)
			},
		)
		go i.Crawl(
			func()(*entity.Stat, error){
				return i.StatRepository.GetLastStatByChainIDAndType(chainID, entity.TypeRemove)
			}, 
			func(lastStat *entity.Stat)([]*entity.Transaction, uint64, error){
				return i.RemoveLiquidCrawl.Crawl(chainID, infoApi, lastStat.Timestamp)
			},
		)
	}
	select { }
}
