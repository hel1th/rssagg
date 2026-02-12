package domain

import (
	"github.com/hel1th/rssagg/internal/database"
)

func MapUserFromDB(dbUser database.User) *User {
	return &User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

func MapFeedFromDB(dbFeed database.Feed) *Feed {
	feed := &Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		URL:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}

	if dbFeed.LastFetchedAt.Valid {
		feed.LastFetchedAt = &dbFeed.LastFetchedAt.Time
	}

	return feed
}

func MapFeedsFromDB(dbFeeds []database.Feed) []*Feed {
	feeds := make([]*Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feeds[i] = MapFeedFromDB(dbFeed)
	}
	return feeds
}

func MapFeedFollowFromDB(dbFeedFollow database.FeedFollow) *FeedFollow {
	return &FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func MapFeedFollowsFromDB(dbFeedFollows []database.FeedFollow) []*FeedFollow {
	feedFollows := make([]*FeedFollow, len(dbFeedFollows))
	for i, dbFeedFollow := range dbFeedFollows {
		feedFollows[i] = MapFeedFollowFromDB(dbFeedFollow)
	}
	return feedFollows
}

func MapPostFromDB(dbPost database.Post) *Post {
	post := &Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		PublishedAt: dbPost.PublishedAt,
		URL:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}

	if dbPost.Description.Valid {
		post.Description = &dbPost.Description.String
	}

	return post
}

func MapPostsFromDB(dbPosts []database.Post) []*Post {
	posts := make([]*Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = MapPostFromDB(dbPost)
	}
	return posts
}
