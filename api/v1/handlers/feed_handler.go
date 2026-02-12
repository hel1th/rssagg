package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/api/v1/dto"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/service"
)

type FeedHandler struct {
	feedService service.FeedService
}

func NewFeedHandler(feedService service.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

func (h *FeedHandler) CreateFeed(w http.ResponseWriter, r *http.Request, user *domain.User) {
	var req dto.CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := h.feedService.CreateFeed(r.Context(), req.Name, req.URL, user.ID)
	if err != nil {
		switch err {
		case domain.ErrInvalidFeedName:
			respondWithError(w, http.StatusBadRequest, "Invalid feed name")
		case domain.ErrInvalidFeedURL:
			respondWithError(w, http.StatusBadRequest, "Invalid feed URL")
		case domain.ErrDuplicateFeed:
			respondWithError(w, http.StatusConflict, "Feed already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create feed: %v", err))
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FeedToResponse(feed))
}

func (h *FeedHandler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.feedService.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get feeds: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dto.FeedsToResponse(feeds))
}

func (h *FeedHandler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
	feedIDStr := r.URL.Query().Get("id")
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

	respondWithJSON(w, http.StatusOK, dto.FeedToResponse(feed))
}
