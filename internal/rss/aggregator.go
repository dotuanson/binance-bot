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
				location, _ := time.LoadLocation("Asia/Bangkok")
				currentTimeUTC7 := item.PublishedParsed.In(location)
				currentTimeUTC7Format := currentTimeUTC7.Format("2006/01/02 - 15:04")
				text := fmt.Sprintf("[%v] %v\nLink: %v", currentTimeUTC7Format, item.Title, item.Link)
				log.Print(text)
				textCh <- text
			}
		}(item)
	}
	time.Sleep(time.Hour * 4)
}
