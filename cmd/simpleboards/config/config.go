package config

import (
	"github.com/spf13/viper"
)

const (
	httpAddr     = "ADDR"
	ddbTablename = "DYNAMODB_TABLE_NAME"
	redisAddr    = "REDIS_ADDR"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("SLBD")
	viper.SetConfigName("simpleboards")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()
	}

	// set defaults
	viper.SetDefault(httpAddr, ":8808")
	viper.SetDefault(redisAddr, "localhost:6379")
	viper.SetDefault(ddbTablename, "sgs-gbl-dev-leaderboards")
}

// GetAddr returns the http server addresss
func GetAddr() string {
	return viper.GetString(httpAddr)
}

// GetDynamoDBTableName returns the http server addresss
func GetDynamoDBTableName() string {
	return viper.GetString(ddbTablename)
}

func GetRedisAddr() string {
	return viper.GetString(redisAddr)
}
