package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hel1th/rssagg/api/v1/dto"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Can't create user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dto.UserToResponse(user))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, user *domain.User) {
	respondWithJSON(w, http.StatusOK, dto.UserToResponse(user))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}
