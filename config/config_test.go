package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	routeTests := []struct {
		name                 string
		environmentVariables map[string]string
		expectedError        error
		output               *Config
	}{
		{
			name: "error without httpPort",
			environmentVariables: map[string]string{
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "5234",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", httpPort),
		},
		{
			name: "error without postgreSQLHost",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "port",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLHost),
		},
		{
			name: "error without postgreSQLUser",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "port",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLUser),
		},
		{
			name: "error without postgreSQLPassword",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "port",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLPassword),
		},
		{
			name: "error without postgreSQLDB",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLPort:                     "port",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLDB),
		},
		{
			name: "error without postgreSQLPort",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLPort),
		},
		{
			name: "error without postgreSQLSubscriptionsTableSchema",
			environmentVariables: map[string]string{
				httpPort:                         "8080",
				postgreSQLHost:                   "host",
				postgreSQLUser:                   "user",
				postgreSQLPassword:               "pass",
				postgreSQLDB:                     "db",
				postgreSQLPort:                   "port",
				postgreSQLSubscriptionsTableName: "subscriptions",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLSubscriptionsTableSchema),
		},
		{
			name: "error without postgreSQLSubscriptionsTableName",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "port",
				postgreSQLSubscriptionsTableSchema: "public",
			},
			expectedError: fmt.Errorf("missing mandatory environment variable: %s", postgreSQLSubscriptionsTableName),
		},
		{
			name: "no error",
			environmentVariables: map[string]string{
				httpPort:                           "8080",
				postgreSQLHost:                     "host",
				postgreSQLUser:                     "user",
				postgreSQLPassword:                 "pass",
				postgreSQLDB:                       "db",
				postgreSQLPort:                     "5234",
				postgreSQLSubscriptionsTableSchema: "public",
				postgreSQLSubscriptionsTableName:   "subscriptions",
			},
			output: &Config{
				HTTPPort: "8080",
				PostgreSQLConfig: PostgreSQLConfig{
					Host:     "host",
					User:     "user",
					Password: "pass",
					DB:       "db",
					Port:     5234,
					SubscriptionsTable: PostgreSQLTableConfig{
						Schema: "public",
						Name:   "subscriptions",
					},
				},
			},
		},
	}

	for _, tt := range routeTests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Unsetenv(httpPort)
			_ = os.Unsetenv(postgreSQLHost)
			_ = os.Unsetenv(postgreSQLUser)
			_ = os.Unsetenv(postgreSQLPassword)
			_ = os.Unsetenv(postgreSQLDB)
			_ = os.Unsetenv(postgreSQLPort)
			_ = os.Unsetenv(postgreSQLSubscriptionsTableSchema)
			_ = os.Unsetenv(postgreSQLSubscriptionsTableName)

			for k, v := range tt.environmentVariables {
				_ = os.Setenv(k, v)
			}

			c, err := InitConfig()
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.output, c)
		})
	}
}
