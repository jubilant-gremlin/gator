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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("ERROR FETCHING REQUEST")
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR FETCHING RESPONSE")
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR READING RESPONSE BODY")
		return &RSSFeed{}, err
	}
	Feed := RSSFeed{}
	err = xml.Unmarshal(dat, &Feed)
	if err != nil {
		fmt.Println("ERROR UNMARSHALING DATA")
		return &RSSFeed{}, err
	}
	Feed.Channel.Title = html.UnescapeString(Feed.Channel.Title)
	Feed.Channel.Description = html.UnescapeString(Feed.Channel.Description)
	for i := range Feed.Channel.Item {
		Feed.Channel.Item[i].Title = html.UnescapeString(Feed.Channel.Item[i].Title)
		Feed.Channel.Item[i].Description = html.UnescapeString(Feed.Channel.Item[i].Description)
	}
	return &Feed, nil
}
