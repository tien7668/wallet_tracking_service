package main

import (
	"fmt"
	"kyberswap_user_monitor/internal/pkg/builder"
	"kyberswap_user_monitor/internal/pkg/state"
	"log"
	"os"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func loadConfig(configFile string) (*state.Cfg, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := &state.Cfg{}
	defaults.SetDefaults(c)
	viper.SetDefault("env", os.Getenv("ENV"))
	viper.SetDefault("database", map[string]interface{}{
		"host":     os.Getenv("DATABASE_HOST"),
		"port":     os.Getenv("DATABASE_PORT"),
		"dBName":   os.Getenv("DATABASE_DBNAME"),
		"user":     os.Getenv("DATABASE_USER"),
		"password": os.Getenv("DATABASE_PASSWORD"),
	})
	fmt.Println(viper.GetString("ENV"))
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	fmt.Printf("config: %+v\n", c)
	return state.GetConfig(c)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "internal/pkg/config/default.yaml",
				Usage:   "Configuration file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "migrate",
				Aliases: []string{},
				Usage:   "migrate database schema",
				Action: func(c *cli.Context) error {
					_, err := loadConfig(c.String("config"))
					if err != nil {
						return err
					}
					cmd, err := builder.NewMigrateBuilder()
					if err != nil {
						return err
					}
					return cmd.Run()
				},
			},
			{
				Name:    "fetcher",
				Aliases: []string{},
				Usage:   "aggregate stats data",
				Action: func(c *cli.Context) error {
					config, err := loadConfig(c.String("config"))
					if err != nil {
						return err
					}
					cmd, err := builder.NewFetcherBuilder(config)
					if err != nil {
						return err
					}
					return cmd.Run()
				},
			},
			{
				Name:    "api",
				Aliases: []string{},
				Usage:   "get stats like trading volume",
				Action: func(c *cli.Context) error {
					config, err := loadConfig(c.String("config"))
					if err != nil {
						return err
					}
					
					cmd, err := builder.NewApiBuilder(config)
					if err != nil {
						return err
					}
					return cmd.Run()
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}