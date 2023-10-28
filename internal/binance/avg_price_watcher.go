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
	QueueLength int = 10
)

func WatchAvgPrice(ctx context.Context, client *binanceConnector.Client, textCh chan<- string, errCh chan<- error, coins []string) {
	prices := make([][]float64, len(coins))
	for i := 0; i < len(coins); i++ {
		prices[i] = make([]float64, 0)
	}
	for {
		for idx, coin := range coins {
			avgPrice, err := client.NewAvgPriceService().Symbol(coin).Do(ctx)
			if err != nil {
				errCh <- err
			}
			price, err := strconv.ParseFloat(avgPrice.Price, 64)
			if err != nil {
				errCh <- err
			}
			prices[idx] = append(prices[idx], price)
			if len(prices[idx]) >= QueueLength {
				diff := prices[idx][QueueLength-1] - prices[idx][0]
				percent := diff / prices[idx][0] * 100
				log.Printf("[%s] CurPrice=%f : increase %.2f%%\n", coin, prices[idx][QueueLength-1], percent)
				if (percent - 1) >= 0 {
					textCh <- fmt.Sprintf(`🚀 [%s] is bullish in 5 minutes, increase 1%%
							=> Current price: %fUSDT\n`, coin, prices[idx][QueueLength-1])
				} else if (-percent - 1) >= 0 {
					textCh <- fmt.Sprintf(`🔥 [%s] is bearish in 5 minutes, decrease 1%%
							=> Current price: %fUSDT\n`, coin, prices[idx][QueueLength-1])
				}
				prices[idx] = prices[idx][1:]
			}
		}
		log.Println()
		time.Sleep(time.Second * 30)
	}
}
