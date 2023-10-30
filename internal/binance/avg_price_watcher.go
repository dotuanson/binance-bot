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
	tick      = time.Second * 10
)

func decreaseOneSecondWatchDogTimer(timer *int64) {
	*timer -= 1
	if *timer < 0 {
		*timer = 0
	}
}

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coin string) {
	var (
		watchdogTimerOnePercentThreshold int64   = 0
		watchdogTimerTwoPercentThreshold int64   = 0
		threshold                        float64 = 1.0
	)
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
		log.Printf("Coin: %s, "+
			"Percent: %.2f %%, "+
			"Close Price: %f, "+
			"Open Price: %f, "+
			"Threshold: %f", coin, percent, closePrice, openPrice, threshold)
		if decreaseOneSecondWatchDogTimer(&watchdogTimerOnePercentThreshold); math.Abs(percent) > 1.0 && watchdogTimerOnePercentThreshold <= 0 {
			watchdogTimerOnePercentThreshold = 5
			textCh <- fmt.Sprintf("%s has just modified %.2f%% in 5m, "+
				"current price: %f USDT\n", coin, percent, closePrice)
		}
		if decreaseOneSecondWatchDogTimer(&watchdogTimerTwoPercentThreshold); math.Abs(percent) > 2.0 && watchdogTimerTwoPercentThreshold <= 0 {
			watchdogTimerTwoPercentThreshold = 10
			if percent >= 0 {
				textCh <- fmt.Sprintf("*🚀 %s is having a bull-run in 5m!*", coin)
			} else {
				textCh <- fmt.Sprintf("*🔥 %s is having a bear-run in 5m!*", coin)
			}
		}
		time.Sleep(tick)
	}
}
