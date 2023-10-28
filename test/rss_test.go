package test

import (
	"github.com/dotuanson/binance-bot/internal/rss"
	"testing"
)

func TestRSS(t *testing.T) {
	textCh := make(chan string, 10)
	rss.FeedRSS("https://cointelegraph.com/rss/tag/bitcoin", textCh)
}
