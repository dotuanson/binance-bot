package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"math"
	"strconv"
	"time"
)

const (
	unitPrice = "USDT"
	tick      = time.Second * 10
)

func decreaseOneSecondWatchDogTimer(timer *int64) {
	*timer -= 1
	if *timer < 0 {
		*timer = 0
	}
}

func WatchSpikePrice(ctx context.Context, client *binanceConnector.Client, threshold float64, textCh chan<- string, errCh chan<- error, coin string) {
	var watchdogTimer int64 = 0
	for {
		kLines, err := client.NewKlinesService().Symbol(coin + unitPrice).
			Interval("1m").Limit(6).
			Do(ctx)
		if err != nil {
			errCh <- err
		}
		if len(kLines) != 6 {
			continue
		}
		lastKLine := kLines[len(kLines)-1]
		beginKLine := kLines[0]
		openPrice, err := strconv.ParseFloat(beginKLine.Open, 64)
		if err != nil {
			errCh <- err
		}
		closePrice, err := strconv.ParseFloat(lastKLine.Close, 64)
		if err != nil {
			errCh <- err
		}
		diff := closePrice - openPrice
		percent := diff / openPrice * 100
		decreaseOneSecondWatchDogTimer(&watchdogTimer)
		if math.Abs(percent) > threshold && watchdogTimer <= 0 {
			watchdogTimer = 6
			if percent >= 0 {
				textCh <- fmt.Sprintf("🚀 %s has just modified %.2f%% in 5m, "+
					"current price: %f USDT\n", coin, percent, closePrice)
			} else {
				textCh <- fmt.Sprintf("🔥 %s has just modified %.2f%% in 5m, "+
					"current price: %f USDT\n", coin, percent, closePrice)
			}
		}
		time.Sleep(tick)
	}
}
