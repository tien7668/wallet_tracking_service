package builder

import (
	"kyberswap_user_monitor/internal/pkg/application"
	"kyberswap_user_monitor/internal/pkg/config"
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/infrastructure/persistence"
	"kyberswap_user_monitor/pkg/context"
	"net/http"

	"kyberswap_user_monitor/internal/pkg/state"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type apiBuilder struct {
	cfg *state.Cfg
	server *gin.Engine
}

func NewApiBuilder(cfg *state.Cfg) (state.IRunner, error) {
	_, err := state.NewDB(&cfg.Database)

	if err != nil {
		return nil, err
	}
	
	application.InitTransactionUsecase(
		application.TransactionUsecase{
			TransactionRepository: persistence.GetTransactionRepoImpl(),
		},
	)

	server, err := NewServer(cfg.Http)

	if err != nil {
		return nil, err
	}

	return &apiBuilder{cfg: cfg, server: server}, nil
}

func NewServer(config *config.Http) (*gin.Engine, error) {
	gin.SetMode(config.Mode)
	server := gin.Default()
	setCORS(server)
	server.GET("/ping", func(c *gin.Context) { c.AbortWithStatus(http.StatusOK) })
	router := server.Group(config.Prefix)
	router.GET("tx", PrefixWithApi(GetTransactionByHash))
	router.GET("txs", PrefixWithApi(GetTransactionsByUserAddress))
	return server, nil
}

func (f *apiBuilder) Run() error {
	return f.server.Run(f.cfg.Http.BindAddress)
}

func PrefixWithApi(f func(context.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := context.New(c).WithLogPrefix("api")
		f(ctx)
	}
}


func setCORS(engine *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowMethods(http.MethodOptions)
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("x-request-id")
	corsConfig.AddAllowHeaders("X-Request-Id")
	engine.Use(cors.New(corsConfig))
}


func GetTransactionsByUserAddress(ctx context.Context) {
	params := struct{
		User string `form:"user"`
	}{}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Errorf("failed to bind params, err: %v", err)
		ctx.AbortWith400(err.Error())
		return
	}
	if params.User == "" {
		ctx.Errorf("failed to bind params, err: Invalid tx hash")
		return
	}
	ctx.Infof("params: %+v", params.User)
	txs, err := application.GetTransactionUsecase().GetByUser(params.User)
	if err != nil {
		ctx.Errorf("failed to get db for transaction, err: %+v", err)
	}

	var txsJSON []entity.TransactionJSON
	for _, tx := range txs {
		txsJSON = append(txsJSON, tx.ToJSON())	
	} 
	ctx.RespondWith(200, "Success", txsJSON)
}

func GetTransactionByHash(ctx context.Context) {
	params := struct{
		TxHash string `form:"tx_hash"`
	}{}
	// ctx.BindHandler(params)
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.Errorf("failed to bind params, err: %v", err)
		ctx.AbortWith400(err.Error())
		return
	}
	if params.TxHash == "" {
		ctx.Errorf("failed to bind params, err: Invalid tx hash")
		return
	}
	
	ctx.Infof("params: %+v", params.TxHash)
	tx, err := application.GetTransactionUsecase().GetByTxHash(params.TxHash)
	ctx.Infof("tx: %+v", tx)
	if err != nil {
		ctx.Errorf("failed to get db for transaction, err: %+v", err)
	}
	ctx.RespondWith(200, "Success", tx)
}

// func SaveTransactionByHash(ctx context.Context) {
// 	params := struct{
// 		TxHash string `form:"tx_hash"`
// 	}{}
// 	// ctx.BindHandler(params)
// 	if err := ctx.ShouldBindQuery(&params); err != nil {
// 		ctx.Errorf("failed to bind params, err: %v", err)
// 		ctx.AbortWith400(err.Error())
// 		return
// 	}
// 	if params.TxHash == "" {
// 		ctx.Errorf("failed to bind params, err: Invalid tx hash")
// 		return
// 	}
// 	err := application.GetTransactionUsecase().TransactionRepository.Save(&entity.Transaction{Id:2, Tx: params.TxHash} )
// 	if err != nil {
// 		ctx.Errorf("failed to get db for transaction, err: %+v", err)
// 	}
// 	ctx.RespondWith(200, "Success", nil)
// }

