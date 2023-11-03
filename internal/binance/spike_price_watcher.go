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
	var (
		watchdogTimerFirstThreshold  int64 = 0
		watchdogTimerSecondThreshold int64 = 0
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
		//log.Printf("Coin: %s, "+
		//	"Percent: %.2f %%, "+
		//	"Close Price: %f, "+
		//	"Open Price: %f, "+
		//	"Threshold: %f", coin, percent, closePrice, openPrice, threshold)
		decreaseOneSecondWatchDogTimer(&watchdogTimerFirstThreshold)
		decreaseOneSecondWatchDogTimer(&watchdogTimerSecondThreshold)
		if math.Abs(percent) > threshold && watchdogTimerFirstThreshold <= 0 {
			watchdogTimerFirstThreshold = 6
			textCh <- fmt.Sprintf("%s has just modified %.2f%% in 5m, "+
				"current price: %f USDT\n", coin, percent, closePrice)
			if math.Abs(percent) > threshold+1.0 && watchdogTimerSecondThreshold <= 0 {
				watchdogTimerSecondThreshold = 6
				if percent >= 0 {
					textCh <- fmt.Sprintf("*🚀 %s is having a bull-run in 5m!*", coin)
				} else {
					textCh <- fmt.Sprintf("*🔥 %s is having a bear-run in 5m!*", coin)
				}
			}
		}
		time.Sleep(tick)
	}
}
