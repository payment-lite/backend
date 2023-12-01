package helpers

import "github.com/spf13/viper"

func LoadENV() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error reading config file")
	}
}
