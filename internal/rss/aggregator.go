package rss

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"sync"
	"time"
)

func getFeed(feedUrl string) *gofeed.Feed {
	fp := gofeed.NewParser()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	return feed
}

func FeedRSS(feedUrls []string, textCh chan<- string) {
	type itemTitleEntity struct {
		title     map[string]bool
		itemMutex sync.RWMutex
	}
	itemTitle := itemTitleEntity{title: make(map[string]bool)}
	for {
		currTime := time.Now()
		for _, feedUrl := range feedUrls {
			feed := getFeed(feedUrl)
			for _, item := range feed.Items {
				go func(item *gofeed.Item) {
					duration := currTime.Sub(*item.PublishedParsed)
					if duration < 12*time.Hour {
						location, _ := time.LoadLocation("Asia/Bangkok")
						currentTimeUTC7 := item.PublishedParsed.In(location)
						currentTimeUTC7Format := currentTimeUTC7.Format("2006/01/02 - 15:04")
						itemTitle.itemMutex.Lock()
						if !itemTitle.title[item.Title] {
							itemTitle.title[item.Title] = true
							text := fmt.Sprintf("[%v] %v\nLink: %v", currentTimeUTC7Format, item.Title, item.Link)
							log.Print(text)
							textCh <- text
						}
						itemTitle.itemMutex.Unlock()
					}
				}(item)
			}
		}
		time.Sleep(time.Minute * 30)
	}
}
