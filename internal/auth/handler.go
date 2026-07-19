package auth

import "net/http"

type Handler struct {
	svc *Service
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", h.login)
	mux.HandleFunc("GET /auth/github/callback", h.callback)
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	state := randomToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300,
	})
	http.Redirect(w, r, h.svc.AuthCodeURL(state), http.StatusFound)
}

func (h *Handler) callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	sess, err := h.svc.HandleCallback(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "auth failed", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		HttpOnly: true,
		Expires:  sess.ExpiresAt,
		Path:     "/",
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
