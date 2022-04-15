package kyberswap

import (
	"fmt"
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/state"
	"github.com/machinebox/graphql"
	"kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl"
)
type AddLiquidCrawlImpl struct{
	GraphClients  map[uint]*graphql.Client
}

const addLiquidQuery = `
	query {
		mints(
			skip: %d,
			first: %d,
			where: {
				timestamp_gt: %d
			},
			orderBy: timestamp
			orderDirection: asc
		){
			id
			pair {
				id
			}
			pool {
				token0 {
					symbol
					decimals
				}
				token1 {
					symbol
					decimals
				}
			}
			amount0
			amount1
			sender
			transaction {
				id
				blockNumber
				timestamp
			}
		}
	}
`

func (c *AddLiquidCrawlImpl) Crawl(chainID uint, tokenInfoApi string, start_time uint64) ([]*entity.Transaction, uint64, error){
	skip := 0
	var transactions []*entity.Transaction
	ctx := state.GetContext()
	lastTimestamp := uint64(0)
	for {
		query := fmt.Sprintf(addLiquidQuery, skip, tx_crawl.GraphFirstLimit, start_time)
		req := graphql.NewRequest(query)
		var resp struct {
			Data []KyberswapLiquidResp `json:"mints"`
		}
		if err := c.GraphClients[uint(chainID)].Run(ctx, req, &resp); err != nil {
			ctx.Errorf("failed to query subgraph, err: %v", err)
			return nil, 0, err
		}

		for i := range resp.Data {
			ex, err := resp.Data[i].ToValidatedResp()
			if err != nil {
				ctx.Errorf("failed to parse kyberswap router response, err: %v", err)
				continue
			}

			if lastTimestamp < ex.Timestamp {
				lastTimestamp = ex.Timestamp
			} 
			transactions = append(transactions, &entity.Transaction{
				ChainID:     uint(chainID),
				Timestamp:   ex.Timestamp,
				Pair:        ex.Pair,
				Extra:       ex.Extra,
				UserAddress: ex.UserAddress,
				BlockNumber: ex.BlockNumber,
				Tx:          ex.Tx,
				LastSync:    ex.Timestamp,
				TxType:		 entity.TxTypeAdd,
			})
		}
		if len(resp.Data) < tx_crawl.GraphFirstLimit {
			ctx.Infoln("no more router exchanges, stop crawling")
			break
		}

		skip += len(resp.Data)
		if skip > tx_crawl.GraphSkipLimit {
			ctx.Infoln("hit skip limit, continue in next cycle")
			ctx.Infoln("data length ex", skip)
			ctx.Infoln("data length transactions", len(transactions))
			break
		}
	}
	return transactions, lastTimestamp, nil
}