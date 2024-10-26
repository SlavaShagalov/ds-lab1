package config

import (
	"github.com/spf13/viper"
)

// Postgres

func SetDefaultPostgresConfig() {
	viper.SetDefault(PostgresHost, "db")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "trello_db")
	viper.SetDefault(PostgresUser, "moderator")
	viper.SetDefault(PostgresPassword, "2222")
	viper.SetDefault(PostgresSSLMode, "disable")
}

func SetTestPostgresConfig() {
	viper.SetDefault(PostgresHost, "test-db")
	viper.SetDefault(PostgresPort, 5432)
	viper.SetDefault(PostgresDB, "trello_db")
	viper.SetDefault(PostgresUser, "moderator")
	viper.SetDefault(PostgresPassword, "2222")
	viper.SetDefault(PostgresSSLMode, "disable")
}
