package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hel1th/rssagg/internal/domain"
)

type Fetcher interface {
	Fetch(url string) (*domain.RSSFeedData, error)
}

type httpFetcher struct {
	client http.Client
}

func NewFetcher() Fetcher {
	return &httpFetcher{
		client: http.Client{Timeout: 10 * time.Second},
	}
}

func (h *httpFetcher) Fetch(url string) (*domain.RSSFeedData, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rssFeed feedXML
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	return xmlToDomain(rssFeed), nil
}

type feedXML struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []itemXML `xml:"item"`
	} `xml:"channel"`
}

type itemXML struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func xmlToDomain(xmlFeed feedXML) *domain.RSSFeedData {
	items := make([]domain.RSSItemData, len(xmlFeed.Channel.Item))
	for i, item := range xmlFeed.Channel.Item {
		items[i] = domain.RSSItemData{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubDate:     item.PubDate,
		}
	}

	return &domain.RSSFeedData{
		Title:       xmlFeed.Channel.Title,
		Description: xmlFeed.Channel.Description,
		Link:        xmlFeed.Channel.Link,
		Items:       items,
	}
}
