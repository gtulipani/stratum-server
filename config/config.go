package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// PostgreSQLTableConfig represents the specific PostgreSQLtable config
type PostgreSQLTableConfig struct {
	Schema string
	Name   string
}

// PostgreSQLConfig represents the specific PostgreSQL config.
type PostgreSQLConfig struct {
	Host               string
	User               string
	Password           string
	DB                 string
	Port               int64
	SubscriptionsTable PostgreSQLTableConfig
}

// Config represents main config.
type Config struct {
	HTTPPort string
	PostgreSQLConfig
}

// InitConfig: loads required configuration
func InitConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	c := Config{
		HTTPPort: v.GetString(httpPort),
		PostgreSQLConfig: PostgreSQLConfig{
			Host:     v.GetString(postgreSQLHost),
			User:     v.GetString(postgreSQLUser),
			Password: v.GetString(postgreSQLPassword),
			DB:       v.GetString(postgreSQLDB),
			Port:     v.GetInt64(postgreSQLPort),
			SubscriptionsTable: PostgreSQLTableConfig{
				Schema: v.GetString(postgreSQLSubscriptionsTableSchema),
				Name:   v.GetString(postgreSQLSubscriptionsTableName),
			},
		},
	}

	if err := validateConfig(v); err != nil {
		return nil, err
	}

	return &c, nil
}

func validateConfig(viper *viper.Viper) error {
	mandatoryVariables := []string{
		httpPort,
		postgreSQLHost,
		postgreSQLUser,
		postgreSQLPassword,
		postgreSQLDB,
		postgreSQLPort,
		postgreSQLSubscriptionsTableSchema,
		postgreSQLSubscriptionsTableName,
	}

	for _, v := range mandatoryVariables {
		if viper.Get(v) == nil {
			return fmt.Errorf("missing mandatory environment variable: %s", v)
		}
	}

	return nil
}
