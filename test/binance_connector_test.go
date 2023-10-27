package test

import (
	"context"
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestHistoricalTradeLookup(t *testing.T) {
	historicalTradeLookup, err := testClient.NewHistoricalTradeLookupService().Symbol("NEARUSDT").Do(context.Background())
	log.Println(binanceConnector.PrettyPrint(historicalTradeLookup))
	require.NoError(t, err)
}

func TestKlines(t *testing.T) {
	klines, err := testClient.NewKlinesService().Symbol("NEARUSDT").Interval("1m").Do(context.Background())
	require.NoError(t, err)
	log.Println(binanceConnector.PrettyPrint(klines))
}

func TestNewAvgPrice(t *testing.T) {
	_, err := testClient.NewAvgPriceService().Symbol("BTCUSDT").Do(context.Background())
	require.NoError(t, err)
}
