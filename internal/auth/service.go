package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	oagithub "golang.org/x/oauth2/github"
)

type UserCreator interface {
	FindOrCreateByGitHubID(ctx context.Context, githubID string, login string, email string) (int64, error)
}

type githubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type githubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type SessionsStore interface {
	Create(ctx context.Context, s Session) error
	Get(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
}

type Service struct {
	oauthCfg *oauth2.Config
	sessions SessionsStore
	users    UserCreator
}

func NewService(clientID, clientSecret string, sessions *SessionRepository, users UserCreator) *Service {
	return &Service{
		oauthCfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     oagithub.Endpoint,
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			Scopes:       []string{"read:user", "user:email"},
		},
		sessions: sessions,
		users:    users,
	}
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func (s *Service) AuthCodeURL(state string) string {
	return s.oauthCfg.AuthCodeURL(state)
}

func (s *Service) HandleCallback(ctx context.Context, code string) (Session, error) {
	token, err := s.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return Session{}, err
	}

	client := s.oauthCfg.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return Session{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("closing response body: %v", err)
		}
	}()
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return Session{}, errors.New("unauthorized access to github user")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Session{}, err
	}
	var gu githubUser
	if err := json.Unmarshal(body, &gu); err != nil {
		return Session{}, err
	}

	response, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return Session{}, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("closing response body: %v", err)
		}
	}()
	if response.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return Session{}, errors.New("unauthorized access to github user")
	}
	emailBody, err := io.ReadAll(response.Body)
	if err != nil {
		return Session{}, err
	}
	var emails []githubEmail
	if err := json.Unmarshal(emailBody, &emails); err != nil {
		return Session{}, err
	}

	var primaryEmail string
	for _, e := range emails {
		if e.Primary {
			primaryEmail = e.Email
			break
		}
	}

	userID, err := s.users.FindOrCreateByGitHubID(ctx, strconv.Itoa(gu.ID), gu.Login, primaryEmail)
	if err != nil {
		return Session{}, err
	}

	sess := Session{
		ID:        randomToken(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.sessions.Create(ctx, sess); err != nil {
		return Session{}, err
	}
	return sess, nil
}

var _ = http.Cookie{}
