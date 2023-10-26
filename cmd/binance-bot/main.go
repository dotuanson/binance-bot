package main

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/dotuanson/binance-bot/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
)

func getAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan string, errCh chan error, coins []string) {
	prices := make([][]float64, len(coins))
	for i := 0; i < len(coins); i++ {
		prices[i] = make([]float64, 0)
	}
	for {
		for idx, coin := range coins {
			avgPrice, err := client.NewAvgPriceService().Symbol(coin).Do(ctx)
			if err != nil {
				errCh <- err
			}
			price, err := strconv.ParseFloat(avgPrice.Price, 64)
			if err != nil {
				errCh <- err
			}
			prices[idx] = append(prices[idx], price)
			if len(prices[idx]) >= 2 {
				diff := prices[idx][1] - prices[idx][0]
				percent := diff / prices[idx][0] * 100
				log.Printf("[%s] - CurPrice = %f - %.2f%%", coin, diff, percent)
				if (percent - 2) >= 0 {
					textCh <- fmt.Sprintf(`🚀 [%s] is bull-run in 2 minutes, increase 2%%
							=> Current price: %fUSDT`, coin, prices[idx][1])
				} else if (-percent - 2) >= 0 {
					textCh <- fmt.Sprintf(`🔥 [%s] is bear-run in 2 minutes, decrease 2%%
							=> Current price: %fUSDT`, coin, prices[idx][1])
				}
				prices[idx] = prices[idx][1:]
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func sendTeleAlert(bot *tgbotapi.BotAPI, teleChatID int64, textCh chan string, errCh chan error) {
	for {
		text := tgbotapi.NewMessage(teleChatID, <-textCh)
		_, err := bot.Send(text)
		if err != nil {
			errCh <- err
		}
	}
}

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
	textCh := make(chan string)
	go getAvgPrice(ctx, client, textCh, errCh, config.CoinLIST)
	go sendTeleAlert(bot, config.TeleCHATID, textCh, errCh)
	for {
		err = <-errCh
		log.Fatal(err)
	}
}
