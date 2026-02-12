package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/service"
)

type RSSHandler struct {
	rssService  service.RSSService
	feedService service.FeedService
}

func NewRSSHandler(rssService service.RSSService, feedService service.FeedService) *RSSHandler {
	return &RSSHandler{
		rssService:  rssService,
		feedService: feedService,
	}
}

func (h *RSSHandler) FetchFeed(w http.ResponseWriter, r *http.Request, user *domain.User) {
	feedIDStr := r.URL.Query().Get("feed_id")
	if feedIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Feed ID is required")
		return
	}

	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed ID format")
		return
	}

	feed, err := h.feedService.GetFeedByID(r.Context(), feedID)
	if err != nil {
		if err == domain.ErrFeedNotFound {
			respondWithError(w, http.StatusNotFound, "Feed not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get feed: %v", err))
		}
		return
	}

	newPostCount, err := h.rssService.FetchSingleFeed(r.Context(), *feed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Feed fetched successfully",
		"feed_id":        feedID,
		"new_post_count": newPostCount,
	})
}
