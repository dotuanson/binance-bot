package rss

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"time"
)

func FeedRSS(feedUrl string, textCh chan<- string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	currTime := time.Now()
	for _, item := range feed.Items {
		go func(item *gofeed.Item) {
			duration := currTime.Sub(*item.PublishedParsed)
			if duration < 12*time.Hour {
				text := fmt.Sprintf(`[%v] - Title: %v
					Link: %v`, *item.PublishedParsed, item.Title, item.Link)
				log.Print(text)
				textCh <- text
			}
		}(item)
	}
	time.Sleep(time.Hour * 4)
}
