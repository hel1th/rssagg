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

type FeedFollowHandler struct {
	feedFollowService service.FeedFollowService
}

func NewFeedFollowHandler(feedFollowService service.FeedFollowService) *FeedFollowHandler {
	return &FeedFollowHandler{
		feedFollowService: feedFollowService,
	}
}

func (h *FeedFollowHandler) FollowFeed(w http.ResponseWriter, r *http.Request, user *domain.User) {
	var req dto.CreateFeedFollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := h.feedFollowService.FollowFeed(r.Context(), user.ID, req.FeedID)
	if err != nil {
		switch err {
		case domain.ErrInvalidFeedID:
			respondWithError(w, http.StatusBadRequest, "Invalid feed ID")
		case domain.ErrDuplicateFeedFollow:
			respondWithError(w, http.StatusConflict, "Already following this feed")
		default:
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to follow feed: %v", err))
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.FeedFollowToResponse(feedFollow))
}

func (h *FeedFollowHandler) GetUserFeedFollows(w http.ResponseWriter, r *http.Request, user *domain.User) {
	feedFollows, err := h.feedFollowService.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get feed follows: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dto.FeedFollowsToResponse(feedFollows))
}

func (h *FeedFollowHandler) UnfollowFeed(w http.ResponseWriter, r *http.Request, user *domain.User) {
	feedFollowIDStr := r.URL.Query().Get("id")
	if feedFollowIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Feed follow ID is required")
		return
	}

	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID format")
		return
	}

	err = h.feedFollowService.UnfollowFeed(r.Context(), feedFollowID, user.ID)
	if err != nil {
		switch err {
		case domain.ErrFeedFollowNotFound:
			respondWithError(w, http.StatusNotFound, "Feed follow not found")
		case domain.ErrCannotUnfollowFeed:
			respondWithError(w, http.StatusForbidden, "Cannot unfollow this feed")
		default:
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unfollow feed: %v", err))
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully unfollowed feed"})
}
