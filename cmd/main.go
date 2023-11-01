package main

import (
	"context"
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/dotuanson/binance-bot/internal/binance"
	"github.com/dotuanson/binance-bot/internal/telegram"
	"github.com/dotuanson/binance-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	FeedUrls = []string{
		"https://cointelegraph.com/rss/tag/bitcoin",
		"https://cointelegraph.com/rss/tag/eth",
	}
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
	for _, coin := range config.CoinLIST {
		if coin == "BTC" {
			go binance.WatchAvgPrice(ctx, client, 0.5, textCh, errCh, coin)

		} else {
			go binance.WatchAvgPrice(ctx, client, 1.5, textCh, errCh, coin)
		}
	}
	for _, coin := range config.CoinLIST {
		go binance.WhaleWatcher(ctx, client, textCh, errCh, coin)
	}
	go telegram.SendTeleAlert(bot, config.TeleCHATID, textCh, errCh)
	//go rss.FeedRSS(FeedUrls, textCh)
	for {
		err = <-errCh
		log.Fatal(err)
	}
}
