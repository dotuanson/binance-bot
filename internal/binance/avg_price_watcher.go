package binance

import (
	"context"
	"fmt"
	binanceConnector "github.com/binance/binance-connector-go"
	"strconv"
	"time"
)

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coins []string) {
	for {
		for _, coin := range coins {
			kLines, err := client.NewKlinesService().Symbol(coin).Interval("5m").Do(ctx)
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
			if (percent - 2) >= 0 {
				textCh <- fmt.Sprintf(`🚀 [%s] is bullish in 5 minutes, increase 2%%
							=> Current high price: USDT\n`, coin)
			} else if (-percent - 2) >= 0 {
				textCh <- fmt.Sprintf(`🔥 [%s] is bearish in 5 minutes, decrease 2%%
							=> Current low price: USDT\n`, coin)
			}
		}
		time.Sleep(time.Second * 30)
	}
}
