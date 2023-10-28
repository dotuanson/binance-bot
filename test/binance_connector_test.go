package test

import (
	"context"
	binanceConnector "github.com/binance/binance-connector-go"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestHistoricalTradeLookup(t *testing.T) {
	_, err := testClient.NewHistoricalTradeLookupService().Symbol("NEARUSDT").Do(context.Background())
	require.NoError(t, err)
}

func TestKlines(t *testing.T) {
	res, err := testClient.NewKlinesService().Symbol("NEARUSDT").Interval("5m").Do(context.Background())
	log.Println(binanceConnector.PrettyPrint(res[len(res)-1]))
	require.NoError(t, err)
}

func TestNewAvgPrice(t *testing.T) {
	_, err := testClient.NewAvgPriceService().Symbol("BTCUSDT").Do(context.Background())
	require.NoError(t, err)
}
