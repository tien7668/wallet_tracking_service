package builder

import (
	"kyberswap_user_monitor/internal/pkg/domain/entity"
	"kyberswap_user_monitor/internal/pkg/state"
	"kyberswap_user_monitor/pkg/context"
)

type migrateBuilder struct{}
func NewMigrateBuilder() (state.IRunner, error) {
	cfg, err := state.GetConfig()
	if err != nil {
		return nil, err
	}
	_, err = state.NewDB(&cfg.Database)
	if err != nil {
		return nil, err
	}
	return &migrateBuilder{}, nil
}

func(b *migrateBuilder) Run() error {
	log := context.NewDefault().WithRequestID(false)
	log = log.WithLogPrefix("migrate")
	log.Infoln("Start migrating....")
	db, err := state.GetDB()
	if err != nil {
		return err
	}
	if err := db.Scopes(entity.TransactionTable()).Migrator().DropTable(&entity.Transaction{}); err != nil {
		log.Errorf("failed to drop volumes table, err: %v", err)
	}
	if err := db.Scopes(entity.TransactionTable()).Migrator().DropTable(&entity.Stat{}); err != nil {
		log.Errorf("failed to drop volumes table, err: %v", err)
	}
	if err := db.Scopes(entity.TransactionTable()).AutoMigrate(&entity.Transaction{}); err != nil {
		log.Errorf("failed to migrate volumes table, err: %v", err)
	}
	if err := db.Scopes(entity.StatTable()).AutoMigrate(&entity.Stat{}); err != nil {
		log.Errorf("failed to migrate stats table, err: %v", err)
	}
	log.Infoln("Finish migrate")
	return nil
}