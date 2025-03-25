package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func cleanRSSData(rss *RSSFeed) {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i := 0; i < len(rss.Channel.Item); i++ {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	my_client := &http.Client{}
	my_request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error with establishing request with context: %w", err)
	}
	my_request.Header.Set("myapp", "gator 0.0")
	res, err := my_client.Do(my_request)
	if err != nil {
		return nil, fmt.Errorf("error with getting response: %w", err)
	}
	res_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body into bytes: %w", err)
	}

	my_rss := &RSSFeed{}
	err = xml.Unmarshal(res_bytes, my_rss)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling into rssfeed: %w", err)
	}
	cleanRSSData(my_rss)

	return my_rss, nil
}
