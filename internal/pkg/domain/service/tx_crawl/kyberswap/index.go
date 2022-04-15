package kyberswap
import (
	// "fmt"
	"strconv"
	// "kyberswap_user_monitor/internal/pkg/domain/entity"
	// "kyberswap_user_monitor/internal/pkg/state"
	// "github.com/machinebox/graphql"
	"encoding/json"
	"kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl"
)
// type KyberswapRouterCrawlImpl struct{
// 	GraphClients  map[uint]*graphql.Client
// }

type KyberswapLiquidResp struct {
	Id          string `json:"id"`
	Pair        struct {
					Id	string `json:"id"`
				} `json:"pair"`
	Pool 		struct {
					Token0	struct{
						Symbol		string `json:"symbol"`
						Decimals	string `json:"decimals"`
					} `json:"token0"`
					Token1	struct{
						Symbol		string `json:"symbol"`
						Decimals	string `json:"decimals"`
					} `json:"token1"`
				} `json:"pool"` 
	Amount0		string `json:amount0` 	
	Amount1		string `json:amount1`
	UserAddress string `json:"sender"`
	Tx 			struct {
					Id          string `json:"id"`
					BlockNumber string `json:"blockNumber"`
					Timestamp   string `json:"timestamp"`
				} `json:"transaction"`
}

func (rr *KyberswapLiquidResp) ToValidatedResp() (*tx_crawl.ValidatedResp, error) {
	timestamp, err := strconv.ParseUint(rr.Tx.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}


	blockNumber, err := strconv.ParseUint(rr.Tx.BlockNumber, 10, 64)
	if err != nil {
		return nil, err
	}

	extraData := map[string]interface{}{
		"token0": rr.Pool.Token0.Symbol,
		"token0Decimals": rr.Pool.Token0.Decimals, 
		"token1": rr.Pool.Token1.Symbol, 
		"token1Decimals": rr.Pool.Token1.Decimals, 
		"amount0": rr.Amount0, 
		"amount1": rr.Amount1,
	}
	extraDataJson, err := json.Marshal(extraData)
	if err != nil {
		return nil, err
	}
	return &tx_crawl.ValidatedResp{
		Id:          rr.Id,
		Pair:        rr.Pair.Id,
		UserAddress: rr.UserAddress,
		Timestamp:   timestamp,
		BlockNumber: blockNumber,
		Tx:          rr.Tx.Id,
		Extra: 		 extraDataJson,
	}, nil
}

// const addLiquidQuery = `
// 	query {
// 		mints(
// 			skip: %d,
// 			first: %d,
// 			where: {
// 				timestamp_gt: %d
// 			},
// 			orderBy: timestamp
// 			orderDirection: asc
// 		){
// 			id
// 			pair {
// 				id
// 			}
// 			pool {
// 				token0 {
// 					symbol
// 					decimals
// 				}
// 				token1 {
// 					symbol
// 					decimals
// 				}
// 			}
// 			amount0
// 			amount1
// 			sender
// 			transaction {
// 				id
// 				blockNumber
// 				timestamp
// 			}
// 		}
// 	}
// `

// const removeLiquidQuery = `
// 	query {
// 		burns(
// 			skip: %d,
// 			first: %d,
// 			where: {
// 				timestamp_gt: %d
// 			},
// 			orderBy: timestamp
// 			orderDirection: asc
// 		){
// 			id
// 			pair {
// 				id
// 			}
// 			pool {
// 				token0 {
// 					symbol
// 					decimals
// 				}
// 				token1 {
// 					symbol
// 					decimals
// 				}
// 			}
// 			amount0
// 			amount1
// 			sender
// 			transaction {
// 				id
// 				blockNumber
// 				timestamp
// 			}
// 		}
// 	}
// `


// type CrawlResp struct {
// 	txs 			[]*entity.Transaction
// 	lastTimestamp	uint64
// 	err				error
// }
// func (c *KyberswapRouterCrawlImpl) CrawlAdd(chainID uint, start_time uint64) (CrawlResp) { 
// 	skip := 0
// 	var transactions []*entity.Transaction
// 	ctx := state.GetContext()
// 	lastTimestamp := uint64(0)
// 	for {
// 		query := fmt.Sprintf(addLiquidQuery, skip, graphFirstLimit, start_time)
// 		req := graphql.NewRequest(query)
// 		var resp struct {
// 			Data []KyberswapLiquidResp `json:"mints"`
// 		}
// 		if err := c.GraphClients[uint(chainID)].Run(ctx, req, &resp); err != nil {
// 			ctx.Errorf("failed to query subgraph, err: %v", err)
// 			return CrawlResp{txs: nil, lastTimestamp: 0, err: err}
// 		}
// 		for i := range resp.Data {
// 			ex, err := resp.Data[i].ToValidatedResp()
// 			if err != nil {
// 				ctx.Errorf("failed to parse kyberswap router response, err: %v", err)
// 				continue
// 			}

// 			if lastTimestamp < ex.Timestamp {
// 				lastTimestamp = ex.Timestamp
// 			} 
// 			transactions = append(transactions, &entity.Transaction{
// 				ChainID:     uint(chainID),
// 				Timestamp:   ex.Timestamp,
// 				Pair:        ex.Pair,
// 				Extra:       ex.Extra,
// 				UserAddress: ex.UserAddress,
// 				BlockNumber: ex.BlockNumber,
// 				Tx:          ex.Tx,
// 				LastSync:    ex.Timestamp,
// 				TxType:		 entity.TxTypeAdd,
// 			})
// 		}
// 		if len(resp.Data) < graphFirstLimit {
// 			ctx.Infoln("no more router exchanges, stop crawling")
// 			break
// 		}

// 		skip += len(resp.Data)
// 		if skip > graphSkipLimit {
// 			ctx.Infoln("hit skip limit, continue in next cycle")
// 			ctx.Infoln("data length ex", skip)
// 			ctx.Infoln("data length transactions", len(transactions))
// 			break
// 		}
// 	}
// 	return CrawlResp{txs: transactions, lastTimestamp: lastTimestamp, err: nil}
// }
// func (c *KyberswapRouterCrawlImpl) CrawlRemove(chainID uint, start_time uint64) (CrawlResp) {
// 	skip := 0
// 	var transactions []*entity.Transaction
// 	ctx := state.GetContext()
// 	lastTimestamp := uint64(0)
// 	for {
// 		query := fmt.Sprintf(removeLiquidQuery, skip, graphFirstLimit, start_time)
// 		req := graphql.NewRequest(query)
// 		var resp struct {
// 			Data []KyberswapLiquidResp `json:"burns"`
// 		}
// 		if err := c.GraphClients[uint(chainID)].Run(ctx, req, &resp); err != nil {
// 			ctx.Errorf("failed to query subgraph, err: %v", err)
// 			return CrawlResp{txs: nil, lastTimestamp: 0, err: err}
// 		}
// 		for i := range resp.Data {
// 			ex, err := resp.Data[i].ToValidatedResp()
// 			if err != nil {
// 				ctx.Errorf("failed to parse kyberswap router response, err: %v", err)
// 				continue
// 			}

// 			if lastTimestamp < ex.Timestamp {
// 				lastTimestamp = ex.Timestamp
// 			} 
// 			transactions = append(transactions, &entity.Transaction{
// 				ChainID:     uint(chainID),
// 				Timestamp:   ex.Timestamp,
// 				Pair:        ex.Pair,
// 				Extra:       ex.Extra,
// 				UserAddress: ex.UserAddress,
// 				BlockNumber: ex.BlockNumber,
// 				Tx:          ex.Tx,
// 				LastSync:    ex.Timestamp,
// 				TxType:		 entity.TxTypeRemove,
// 			})
// 		}
// 		if len(resp.Data) < graphFirstLimit {
// 			ctx.Infoln("no more router exchanges, stop crawling")
// 			break
// 		}

// 		skip += len(resp.Data)
// 		if skip > graphSkipLimit {
// 			ctx.Infoln("hit skip limit, continue in next cycle")
// 			ctx.Infoln("data length ex", skip)
// 			ctx.Infoln("data length transactions", len(transactions))
// 			break
// 		}
// 	}
// 	return CrawlResp{txs: transactions, lastTimestamp: lastTimestamp, err: nil}
// }

// func (c *KyberswapRouterCrawlImpl) Crawl(chainID uint, tokenInfoApi string, start_time uint64) ([]*entity.Transaction, uint64, error){
// 	var transactions []*entity.Transaction
// 	chanel := make(chan CrawlResp)
// 	go func(chanel chan CrawlResp){
// 		chanel <- c.CrawlAdd(chainID, start_time)
// 	}(chanel)
// 	go func(chanel chan CrawlResp){
// 		chanel <- c.CrawlRemove(chainID, start_time)
// 	}(chanel)
// 	addResp, removeResp := <-chanel, <-chanel
// 	lastTimestamp := addResp.lastTimestamp
// 	if addResp.err != nil {
// 		return addResp.txs, 0, addResp.err
// 	}
// 	if removeResp.err != nil {
// 		return removeResp.txs, 0, removeResp.err
// 	}
// 	transactions = append(addResp.txs, removeResp.txs...)
// 	if addResp.lastTimestamp < removeResp.lastTimestamp {
// 		lastTimestamp = removeResp.lastTimestamp
// 	}
// 	return transactions, lastTimestamp, nil
// }
