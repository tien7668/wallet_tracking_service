package state

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
type dbConfig struct {
	Host         string `default:"localhost"`
	Port         int    `default:"3306"`
	DBName       string `default:"database"`
	User         string `default:"user"`
	Password     string `default:"password"`
	ConnLifeTime int    `default:"300"`
	ConnTimeOut  int    `default:"30"`
	MaxIdleConns int    `default:"10"`
	MaxOpenConns int    `default:"80"`
	LogLevel     int    `default:"1"`
}

func (c *dbConfig) DNS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ds",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.ConnTimeOut,
	)
}

func NewDB(cfg *dbConfig) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(cfg.DNS()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open connection, err: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.db, err: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnLifeTime) * time.Second)

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database, err: %v", err)
	}
	return GetDB(db)
}
