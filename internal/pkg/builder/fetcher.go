package builder

import (
	"kyberswap_user_monitor/internal/pkg/application"
	tokenService "kyberswap_user_monitor/internal/pkg/domain/service/token_info"
	swapCrawlService "kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl/aggregator"
	liquidCrawlService "kyberswap_user_monitor/internal/pkg/domain/service/tx_crawl/kyberswap"
	"kyberswap_user_monitor/internal/pkg/infrastructure/persistence"
	"kyberswap_user_monitor/internal/pkg/state"
	"kyberswap_user_monitor/pkg/context"

	"github.com/machinebox/graphql"
)

type fetcherBuilder struct {
	cfg *state.Cfg
	transactionUseCase application.TransactionUsecase
}

func NewFetcherBuilder(cfg *state.Cfg) (state.IRunner, error) {
	_, err := state.NewDB(&cfg.Database)

	if err != nil {
		return nil, err
	}
	aggregatorGraphClients := make(map[uint]*graphql.Client, 0)
	kyberswapGraphClients := make(map[uint]*graphql.Client, 0)
	for _, chain := range cfg.Blockchain.Chains {
		aggregatorGraphClients[chain.ChainID] = graphql.NewClient(chain.Subgraph.AggregatorRouter)
		kyberswapGraphClients[chain.ChainID] = graphql.NewClient(chain.Subgraph.KyberswapRouter)
	}

	// graphClient := graphql.NewClient(cfg.Blockchain.Chains[0].Subgraph.AggregatorRouter)
	tokensMap := map[string]tokenService.TokenInfo {}
	application.InitTransactionUsecase(
		application.TransactionUsecase{
			TransactionRepository: persistence.GetTransactionRepoImpl(),
			StatRepository: persistence.GetStatRepoImpl(),
			SwapCrawl: &swapCrawlService.SwapCrawlImpl{ GraphClients: aggregatorGraphClients, Tokens: tokensMap },
			AddLiquidCrawl: &liquidCrawlService.AddLiquidCrawlImpl{ GraphClients: kyberswapGraphClients},
			RemoveLiquidCrawl: &liquidCrawlService.RemoveLiquidCrawlImpl{ GraphClients: kyberswapGraphClients},
			Tokens: tokensMap,
		},
	)
	return &fetcherBuilder{cfg: cfg, transactionUseCase: application.GetTransactionUsecase()}, nil
}

func (b *fetcherBuilder) Run() error {
	state.InitContext(context.NewDefault().WithRequestID(false)) 
	return b.transactionUseCase.CrawAll()
}