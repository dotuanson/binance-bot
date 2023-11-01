package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"strconv"
	"time"
)

func WhaleWatcher(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coin string) {
	const (
		ratioVolume    = 30.0
		numberOfKlines = 3
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
		middlePrice, err := strconv.ParseFloat(kLines[numberOfKlines-2].Close, 64)
		if err != nil {
			errCh <- err
		}
		middleVolume, err := strconv.ParseFloat(kLines[numberOfKlines-2].Volume, 64)
		if err != nil {
			errCh <- err
		}
		firstPrice, err := strconv.ParseFloat(kLines[numberOfKlines-1].Close, 64)
		if err != nil {
			errCh <- err
		}
		firstVolume, err := strconv.ParseFloat(kLines[numberOfKlines-3].Volume, 64)
		if err != nil {
			errCh <- err
		}
		if lastVolume < 1e-9 {
			continue
		}
		//log.Printf("Coin: %s, "+
		//	"Current Price: %f, "+
		//	"Last/ Middle Volume: %f, "+
		//	"Last/ First Volume: %f", coin, lastPrice, lastVolume/middleVolume, lastVolume/firstVolume)
		if lastVolume > 2*ratioVolume*firstVolume {
			diff := lastPrice - firstPrice
			percent := diff / firstPrice * 100
			if percent >= 1 {
				textCh <- fmt.Sprintf("*[x%f] 🚀 %s is having a bull-run with %.2f%%, "+
					"current price: %f USDT*", 2*ratioVolume, coin, percent, lastPrice)
			} else if -percent >= 1 {
				textCh <- fmt.Sprintf("*[x%f] 🔥 %s is having a bull-run with %.2f%%, "+
					"current price: %f USDT*", 2*ratioVolume, coin, percent, lastPrice)
			}
			time.Sleep(time.Second * 30)
		} else if lastVolume > ratioVolume*middleVolume {
			diff := lastPrice - middlePrice
			percent := diff / middlePrice * 100
			if percent >= 1 {
				textCh <- fmt.Sprintf("[x%f] 🚀 %s has just modified %.2f%%, "+
					"current price: %f USDT", ratioVolume, coin, percent, lastPrice)
			} else if -percent >= 1 {
				textCh <- fmt.Sprintf("[x%f] 🔥 %s has just modified %.2f%%, "+
					"current price: %f USDT", ratioVolume, coin, percent, lastPrice)
			}
			time.Sleep(time.Second * 20)
		}
		time.Sleep(time.Second * 1)
	}
}
