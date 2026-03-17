package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/auth"
)

type AuthController struct {
	app   *application.Application
	store sessions.Store
}

func NewAuthController(app *application.Application, store sessions.Store) *AuthController {
	return &AuthController{app, store}
}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	auth.Login().Render(r.Context(), w)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	user, err := c.app.SyncUser(r.Context(), u)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to sync user: %w", err))
		return
	}

	session, _ := c.store.Get(r, middlewares.SessionKey)
	session.Values[middlewares.UserIDKey] = user.ID.Hex()

	redirect := "/"
	next := session.Flashes(middlewares.RedirectKey)
	if len(next) > 0 && next[0].(string) != "" {
		redirect = next[0].(string)
	}

	err = session.Save(r, w)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to save session: %w", err))
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to logout: %w", err))
		return
	}

	session, _ := c.store.Get(r, middlewares.SessionKey)
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)
	err = session.Save(r, w)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to clear session: %w", err))
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (c *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("authentication failed: %w", err))
		return
	}

	user, err := c.app.SyncUser(r.Context(), u)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to sync user: %w", err))
		return
	}

	session, _ := c.store.Get(r, middlewares.SessionKey)
	session.Values[middlewares.UserIDKey] = user.ID.Hex()

	redirect := "/dashboard"
	next := session.Flashes(middlewares.RedirectKey)
	if len(next) > 0 && next[0].(string) != "" {
		redirect = next[0].(string)
	}

	err = session.Save(r, w)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to save session: %w", err))
		return
	}

	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}
