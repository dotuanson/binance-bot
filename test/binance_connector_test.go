package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHistoricalTradeLookup(t *testing.T) {
	_, err := testClient.NewHistoricalTradeLookupService().Symbol("NEARUSDT").Do(context.Background())
	require.NoError(t, err)
}

func TestKlines(t *testing.T) {
	_, err := testClient.NewKlinesService().Symbol("NEARUSDT").Interval("1m").Do(context.Background())
	require.NoError(t, err)
}

func TestNewAvgPrice(t *testing.T) {
	_, err := testClient.NewAvgPriceService().Symbol("BTCUSDT").Do(context.Background())
	require.NoError(t, err)
}
