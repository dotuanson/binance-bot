package test

import (
	"github.com/dotuanson/binance-bot/internal/rss"
	"testing"
	"time"
)

func TestRSS(t *testing.T) {
	textCh := make(chan string, 10)
	go rss.FeedRSS("https://cointelegraph.com/rss/tag/bitcoin", textCh)
	time.Sleep(time.Second * 10)
}
