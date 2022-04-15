package tx_crawl
import (
	"kyberswap_user_monitor/internal/pkg/domain/entity"
)
type CrawlInterface interface {
	Crawl(chainID uint, tokenInfoApi string, start_time uint64) ([]*entity.Transaction, uint64, error)
}

type ValidatedResp struct {
	Id          string
	Pair        string
	UserAddress string
	Timestamp   uint64
	BlockNumber uint64
	Tx          string
	Extra		[]byte
}

const (
	GraphSkipLimit  = 5000
	GraphFirstLimit = 500
)