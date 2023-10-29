package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"log"
	"strconv"
	"time"
)

const (
	unitPrice = "USDT"
)

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coin string) {
	var (
		thresholdUp   float64 = 1.0
		thresholdDown float64 = 1.0
	)
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
			"Threshold Up: %f, "+
			"Threshold Down: %f", coin, percent, closePrice, openPrice, thresholdUp, thresholdDown)
		if (percent - thresholdUp) >= 0 {
			textCh <- fmt.Sprintf("🚀 [%s] is BULLISH (up %.2f%%) in 5m, "+
				"price: %f USDT\n", coin, thresholdUp, closePrice)
			thresholdUp += 0.1
			time.Sleep(time.Second * 25)
		} else if (-percent - thresholdDown) >= 0 {
			textCh <- fmt.Sprintf("🔥 [%s] is BEARISH (down %.2f%%) in 5m, "+
				"price: %f USDT\n", coin, thresholdDown, closePrice)
			thresholdDown += 0.1
			time.Sleep(time.Second * 25)
		} else {
			if thresholdUp -= 0.1; (thresholdUp - 1.0) <= 0.0 {
				thresholdUp = 1.0
			} else {
				time.Sleep(time.Second * 25)
			}
			if thresholdDown -= 0.1; (thresholdDown - 1.0) <= 0.0 {
				thresholdDown = 1.0
			} else {
				time.Sleep(time.Second * 25)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
