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

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coins []string) {
	for {
		for _, coin := range coins {
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
				"Close Price %f, "+
				"Open Price %f", coin, percent, closePrice, openPrice)
			if (percent - 1) >= 0 {
				textCh <- fmt.Sprintf("🚀 [%s] is BULLISH (up 1%%) in 5m, "+
					"price: %f USDT\n", coin, closePrice)
			} else if (-percent - 1) >= 0 {
				textCh <- fmt.Sprintf("🔥 [%s] is BEARISH (down 1%%) in 5m, "+
					"price: %f USDT\n", coin, closePrice)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
