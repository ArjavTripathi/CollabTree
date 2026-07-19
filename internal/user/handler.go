package user

import (
	"SocialMedia/internal/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("GET /users/me", authMW(http.HandlerFunc(h.GetProfile)))
	//mux.Handle("POST /users/create", authMW(http.HandlerFunc(h.CreateUser))) DEPRECATED due to new login handler
	mux.Handle("DELETE /users/me", authMW(http.HandlerFunc(h.DeleteUser)))
	mux.Handle("PUT /users/me", authMW(http.HandlerFunc(h.UpdateProfile)))

}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var input User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, errId := middleware.GetUserIdFromContext(r.Context())
	if errId != nil {
		http.Error(w, error.Error(errId), http.StatusBadRequest)
	}

	err := h.service.UpdateProfile(r.Context(), &input, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "user updated"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, errId := middleware.GetUserIdFromContext(r.Context())
	if errId != nil {
		http.Error(w, error.Error(errId), http.StatusBadRequest)
	}

	u, err := h.service.GetProfile(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateUser(r.Context(), &input)
	if err != nil {
		http.Error(w, "internal error in service layer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]int64{"id": id})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, errId := middleware.GetUserIdFromContext(r.Context())
	if errId != nil {
		http.Error(w, error.Error(errId), http.StatusBadRequest)
	}

	err := h.service.DeleteUser(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "user not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("Deleted id: %d", id)})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
