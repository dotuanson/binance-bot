package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"strconv"
	"time"
)

func WatchHighVolume(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coin string) {
	const (
		ratioVolume    = 30.0
		numberOfKlines = 2
	)
	for {
		kLines, err := client.NewKlinesService().Symbol(coin + unitPrice).
			Interval("1m").Limit(numberOfKlines).
			Do(ctx)
		if err != nil {
			errCh <- err
		}
		if len(kLines) != numberOfKlines {
			continue
		}
		lastPrice, err := strconv.ParseFloat(kLines[numberOfKlines-1].Close, 64)
		if err != nil {
			errCh <- err
		}
		lastVolume, err := strconv.ParseFloat(kLines[numberOfKlines-1].Volume, 64)
		if err != nil {
			errCh <- err
		}
		firstPrice, err := strconv.ParseFloat(kLines[0].Close, 64)
		if err != nil {
			errCh <- err
		}
		firstVolume, err := strconv.ParseFloat(kLines[0].Volume, 64)
		if err != nil {
			errCh <- err
		}
		if lastVolume < 1e-9 {
			continue
		}
		if lastVolume > 2*ratioVolume*firstVolume {
			diff := lastPrice - firstPrice
			percent := diff / firstPrice * 100
			if percent >= 1 {
				textCh <- fmt.Sprintf("*[x%f] 🚀 %s is having a bull-run with %.2f%%, "+
					"current price: %f USDT*", lastVolume/firstVolume, coin, percent, lastPrice)
			} else if -percent >= 1 {
				textCh <- fmt.Sprintf("*[x%f] 🔥 %s is having a bull-run with %.2f%%, "+
					"current price: %f USDT*", lastVolume/firstVolume, coin, percent, lastPrice)
			}
			time.Sleep(time.Second * 60)
		}
		time.Sleep(time.Second * 10)
	}
}
