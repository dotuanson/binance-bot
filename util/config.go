package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	ApiKEY     string `mapstructure:"API_KEY"`
	SecretKEY  string `mapstructure:"SECRET_KEY"`
	BaseURL    string `mapstructure:"BASE_URL"`
	TeleTOKEN  string `mapstructure:"TELEGRAM_TOKEN"`
	TeleCHATID int64
	CoinLIST   []string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	config.TeleCHATID = viper.GetInt64("TELEGRAM_CHATID")
	config.CoinLIST = viper.GetStringSlice("COIN_LIST")
	return
}
