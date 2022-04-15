package state

import (
	"fmt"
	"kyberswap_user_monitor/pkg/context"

	"gorm.io/gorm"
)

var ctx context.Context

var dbInstance *gorm.DB

var configInstance *Cfg

func GetContext() context.Context{
	return ctx
} 

func InitContext(c context.Context) {
	ctx = c
} 

func GetDB(db ...*gorm.DB) (*gorm.DB, error){

	if dbInstance != nil {
		return dbInstance, nil
	}

	if len(db) == 0 {
		return nil, fmt.Errorf("database not init")
	}

	dbInstance = db[0]
	return dbInstance, nil
}

func GetConfig(cfg ...*Cfg) (*Cfg, error) {
	if configInstance != nil {
		return configInstance, nil
	}

	if len(cfg) == 0 {
		return nil, fmt.Errorf("config not init")
	}
	configInstance = cfg[0]
	return configInstance, nil
}
