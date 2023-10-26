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

func getAvgPrice(ctx context.Context, client *binanceConnector.Client, c chan string, errCh chan error) {
	prices := make([]float64, 0)
	for {
		avgPrice, err := client.NewAvgPriceService().Symbol("NEARUSDT").Do(ctx)
		if err != nil {
			errCh <- err
		}
		price, err := strconv.ParseFloat(avgPrice.Price, 64)
		if err != nil {
			errCh <- err
		}
		prices = append(prices, price)
		if len(prices) >= 2 {
			log.Printf("%f - %f - %f", prices[0], prices[1], prices[1]-prices[0])
			if ((prices[1]-prices[0])/prices[0] - 0.02) >= 0 {
				c <- fmt.Sprintf(`
					🚀 %s is bull-run in 2 minutes, increase 2%%\n
				`, "NEAR")
			} else if ((prices[0]-prices[1])/prices[0] - 0.02) >= 0 {
				c <- fmt.Sprintf(`
					🔥 %s is bear-run in 2 minutes, decrease 2%%\n
				`, "NEAR")
			}
			prices = prices[1:]
		}
		time.Sleep(time.Minute * 2)
	}
}

func sendTeleAlert(bot *tgbotapi.BotAPI, teleChatID int64, c chan string, errCh chan error) {
	for {
		text := binanceConnector.PrettyPrint(<-c)
		msg := tgbotapi.NewMessage(teleChatID, text)
		_, err := bot.Send(msg)
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
	go getAvgPrice(ctx, client, textCh, errCh)
	go sendTeleAlert(bot, config.TeleCHATID, textCh, errCh)
	for {
		err = <-errCh
		log.Fatal(err)
	}
}
