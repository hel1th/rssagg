package domain

type RSSFeedData struct {
	Title       string
	Description string
	Link        string
	Items       []RSSItemData
}

type RSSItemData struct {
	Title       string
	Description string
	Link        string
	PubDate     string
}
