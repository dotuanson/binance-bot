package main

import (
	"context"
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/dotuanson/binance-bot/internal/binance"
	"github.com/dotuanson/binance-bot/internal/rss"
	"github.com/dotuanson/binance-bot/internal/telegram"
	"github.com/dotuanson/binance-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	FeedUrl = "https://cointelegraph.com/rss/tag/bitcoin"
)

func main() {
	ctx := context.Background()
	errCh := make(chan error)
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TeleTOKEN)
	if err != nil {
		log.Fatal(err)
	}

	client := binanceConnector.NewClient(config.ApiKEY, config.SecretKEY, config.BaseURL)
	textCh := make(chan string, 10)
	go binance.GetAvgPrice(ctx, client, textCh, errCh, config.CoinLIST)
	go telegram.SendTeleAlert(bot, config.TeleCHATID, textCh, errCh)
	go rss.FeedRSS(FeedUrl, textCh)
	for {
		err = <-errCh
		log.Fatal(err)
	}
}
