package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	unitPrice = "USDT"
)

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coin string) {
	var threshold = 1.0
	for {
		kLines, err := client.NewKlinesService().Symbol(coin + unitPrice).Interval("5m").Do(ctx)
		if err != nil {
			errCh <- err
		}
		lastKLine := kLines[len(kLines)-1]
		openPrice, err := strconv.ParseFloat(lastKLine.Open, 64)
		if err != nil {
			errCh <- err
		}
		closePrice, err := strconv.ParseFloat(lastKLine.Close, 64)
		if err != nil {
			errCh <- err
		}
		diff := closePrice - openPrice
		percent := diff / openPrice * 100
		log.Printf("Coin: %s, "+
			"Percent: %.2f %%, "+
			"Close Price: %f, "+
			"Open Price: %f, "+
			"Threshold: %f", coin, percent, closePrice, openPrice, threshold)
		if math.Pow(percent, 2)-math.Pow(threshold, 2) > 0 {
			if (threshold - 1.0) <= 1e-9 {
				textCh <- fmt.Sprintf("[%s] has just modified %.2f%% in 5m, "+
					"price: %f USDT\n", coin, percent, closePrice)
				threshold = percent
				time.Sleep(time.Second * 25)
			} else {
				if percent >= 0 {
					textCh <- fmt.Sprintf("🚀 [%s] is having a bull-run in 5m!", coin)
				} else {
					textCh <- fmt.Sprintf("🔥 [%s] is having a bear-run in 5m!", coin)
				}
				threshold = 1.0
				time.Sleep(time.Second * 55)
			}
		} else {
			threshold = 1.0
		}
		time.Sleep(time.Second * 5)
	}
}
