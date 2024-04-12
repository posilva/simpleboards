package config

import (
	"github.com/spf13/viper"
)

const (
	httpAddr = "ADDR"
)

func init() {
	viper.SetDefault(httpAddr, ":8808")

	viper.AddConfigPath(".")
	viper.SetEnvPrefix("SLBD")
	viper.SetConfigName("simpleboards")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()
	}
}

// GetAddr returns the http server addresss
func GetAddr() string {
	return viper.GetString(httpAddr)
}
