package test

import (
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/dotuanson/binance-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"testing"
)

var (
	testClient   *binanceConnector.Client
	testTgBotAPI *tgbotapi.BotAPI
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testClient = binanceConnector.NewClient(config.ApiKEY, config.SecretKEY, config.BaseURL)
	testTgBotAPI, err = tgbotapi.NewBotAPI(config.TeleTOKEN)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
