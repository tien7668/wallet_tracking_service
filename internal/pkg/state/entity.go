package state

import (
	"kyberswap_user_monitor/internal/pkg/config"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type IRunner interface {
	Run() error
}

type Cfg struct {
	Env string
	Http *config.Http
	Database    dbConfig
	Blockchain *config.Blockchain
}