package devtools

import (
	"SocialMedia/internal/auth"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SessionCreator interface {
	Create(ctx context.Context, session auth.Session) error
}

type Handler struct {
	sessions SessionCreator
}

func NewHandler(sessions SessionCreator) *Handler {
	return &Handler{sessions: sessions}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /dev/login", h.devLogin)
}

func (h *Handler) devLogin(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("userId")
	if userIDParam == "" {
		userIDParam = "3"
	}
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	token, err := randomToken()
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	session := auth.Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.sessions.Create(r.Context(), session); err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "X-Session-ID",
		Value:    token,
		HttpOnly: true,
		Expires:  session.ExpiresAt,
		Path:     "/",
	})

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
