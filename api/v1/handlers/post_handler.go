package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hel1th/rssagg/api/v1/dto"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) GetPostsForUser(w http.ResponseWriter, r *http.Request, user *domain.User) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	posts, err := h.postService.GetPostsForUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get posts: %v", err))
		return
	}

	postValues := make([]domain.Post, len(posts))
	for i, post := range posts {
		postValues[i] = *post
	}

	respondWithJSON(w, http.StatusOK, dto.PostsToResponse(postValues))
}
