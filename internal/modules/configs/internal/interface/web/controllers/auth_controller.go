package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
)

const (
	SessionName = "session"
	UserIDKey   = "user_id"
	RedirectKey = "redirect"
)

type AuthController struct {
	app   *application.Application
	store sessions.Store
}

func NewAuthController(app *application.Application, store sessions.Store) *AuthController {
	return &AuthController{app, store}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	user, err := c.app.SyncUser(r.Context(), u)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sync user: %v", err), http.StatusInternalServerError)
		return
	}

	session, _ := c.store.Get(r, SessionName)
	session.Values[UserIDKey] = user.ID.Hex()

	redirect := "/"
	next := session.Flashes(RedirectKey)
	if len(next) > 0 && next[0].(string) != "" {
		redirect = next[0].(string)
	}

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	session, _ := c.store.Get(r, SessionName)
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (c *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}

	user, err := c.app.SyncUser(r.Context(), u)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sync user: %v", err), http.StatusInternalServerError)
		return
	}

	session, _ := c.store.Get(r, SessionName)
	session.Values[UserIDKey] = user.ID.Hex()

	redirect := "/dashboard"
	next := session.Flashes(RedirectKey)
	if len(next) > 0 && next[0].(string) != "" {
		redirect = next[0].(string)
	}

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}
