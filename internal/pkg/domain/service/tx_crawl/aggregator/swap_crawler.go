package aggregator

import (
	"fmt"
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/state"
	"kyberswap_user_monitor/internal/pkg/config"
	"kyberswap_user_monitor/internal/pkg/domain/service/token_info"
	"strconv"
	"encoding/json"
	"github.com/machinebox/graphql"
	"strings"
	"sync"
	"kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl"

)

type SwapCrawlImpl struct {
	GraphClients  map[uint]*graphql.Client
	Tokens map[string]token_info.TokenInfo
	mu sync.Mutex 
}

type RouterResp struct {
	Id          	string `json:"id"`
	Pair        	string `json:"pair"`
	TokenIn       	string `json:"tokenIn"`
	AmountIn      	string `json:"amountIn"`
	TokenOut       	string `json:"tokenOut"`
	AmountOut      	string `json:"amountOut"`
	UserAddress 	string `json:"userAddress"`
	Timestamp   	string `json:"time"`
	BlockNumber 	string `json:"blockNumber"`
	Tx          	string `json:"tx"`
}

// type RouterExchange struct {
// 	Id          string
// 	Pair        string
// 	UserAddress string
// 	Timestamp   uint64
// 	BlockNumber uint64
// 	Tx          string
// 	Extra		[]byte
// }


func (rr *RouterResp) ToValidatedResp(chainID uint, tokenInfo map[string]token_info.TokenInfo) (*tx_crawl.ValidatedResp, error) {
	exchange := &entity.Transaction{}
	timestamp, err := strconv.ParseUint(rr.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	exchange.Timestamp = timestamp
	blockNumber, err := strconv.ParseUint(rr.BlockNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	
	idIn := fmt.Sprint(chainID) + "_" + rr.TokenIn
	idOut := fmt.Sprint(chainID) + "_" + rr.TokenOut
	extraInfoIn := map[string]interface{} {}
	extraInfoOut := map[string]interface{} {}
	if (rr.TokenIn == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee") {
		extraInfoIn["tokenInSymbol"] = config.ETH[chainID].Symbol
		extraInfoIn["tokenInDecimals"] = config.ETH[chainID].Decimals
	} else {
		tIn, existed := tokenInfo[idIn]
		if existed {
			extraInfoIn["tokenInSymbol"] = tIn.Symbol
			extraInfoIn["tokenInDecimals"] = tIn.Decimals
		}
	}
	if (rr.TokenOut == "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"){
		extraInfoOut["tokenOutSymbol"] = config.ETH[chainID].Symbol
		extraInfoOut["tokenOutDecimals"] = config.ETH[chainID].Decimals
	}else {
		tOut, existed := tokenInfo[idOut]
		if existed {
			extraInfoOut["tokenOutSymbol"] = tOut.Symbol
			extraInfoOut["tokenOutDecimals"] = tOut.Decimals
		}
	}

	extraData := map[string]interface{}{
		"tokenIn": rr.TokenIn, 
		"amountIn": rr.AmountIn,
		"tokenOut": rr.TokenOut, 
		"amountOut": rr.AmountOut,
	}
	for k, v := range extraInfoIn {
		extraData[k] = v
	}
	for k, v := range extraInfoOut {
		extraData[k] = v
	}

	extraDataJson, err := json.Marshal(extraData)

	if err != nil {
		return nil, err
	}
	return &tx_crawl.ValidatedResp{
		Id:          rr.Id,
		Pair:        rr.Pair,
		UserAddress: rr.UserAddress,
		Timestamp:   timestamp,
		BlockNumber: blockNumber,
		Tx:          rr.Tx,
		Extra: 		 extraDataJson,
	}, nil
} 


const routerLogQuery = `
	query {
		routerSwappeds(
			skip: %d,
			first: %d,
			where: {
				time_gt: %d
			},
			orderBy: time
			orderDirection: asc
		) {
			id
			pair
			tokenIn
			amountIn
			tokenOut
			amountOut
			userAddress
			time
			blockNumber
			tx
		}
	}
`


func (c *SwapCrawlImpl) Crawl(chainID uint, tokenInfoApi string, start_time uint64) ([]*entity.Transaction, uint64, error){
	skip := 0
	// var exchanges []RouterResp
	// tokens := make(map[string]*httputils.TokenInfo)
	var transactions []*entity.Transaction
	ctx := state.GetContext()
	lastTimestamp := uint64(0)
	
	flags := map[string]bool{}
	for {
		query := fmt.Sprintf(routerLogQuery, skip, tx_crawl.GraphFirstLimit, start_time)
		req := graphql.NewRequest(query)
		var resp struct {
			Data []RouterResp `json:"routerSwappeds"`
		}
		if err := c.GraphClients[uint(chainID)].Run(ctx, req, &resp); err != nil {
			ctx.Errorf("failed to query subgraph, err: %v", err)
			return nil, 0, err
		}
		// exchanges = append(exchanges, resp.Data...)
		
		ids := []string{}
		for i := range resp.Data {
			idIn := fmt.Sprint(chainID) + "_" + resp.Data[i].TokenIn
			idOut := fmt.Sprint(chainID) + "_" + resp.Data[i].TokenOut
			
			if _, existed := c.Tokens[idIn] ; !existed && !flags[idIn]{
				ids = append(ids, resp.Data[i].TokenIn)
			}
			if _, existed := c.Tokens[idOut] ; !existed && !flags[idOut]{
				ids = append(ids, resp.Data[i].TokenOut)
			}
			flags[idIn] = true
			flags[idOut] = true
		}
		concatedIDs := strings.Join(ids[:], ",")

		tokens, err := token_info.GetTokenInfo(tokenInfoApi, concatedIDs)
		if err == nil {
			c.mu.Lock()
			for k, v := range tokens {
				id := fmt.Sprint(chainID) + "_" + k
				c.Tokens[id] = *v
			}
			c.mu.Unlock()
		}

		for i := range resp.Data {
			ex, err := resp.Data[i].ToValidatedResp(chainID, c.Tokens)
			if err != nil {
				ctx.Errorf("failed to parse router response, err: %v", err)
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
				TxType:		 entity.TxTypeSwap,
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
