package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	}
}

// LoginHandler redirects user to oauth flow
func LoginHandler(auth *Auth, force bool, offline bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := auth.BuildRedirectURL(force, offline)
		sessionID := redirectURL.SessionID
		oauthFlowSession, err := getOauthFlowSession(r, sessionID)
		handleError(w, err)

		redirectURLStr, err := validateRedirectURL(r.FormValue(oauthFlowRedirectKey))
		handleError(w, err)

		oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURLStr
		err = oauthFlowSession.Save(r, w)
		handleError(w, err)

		http.Redirect(w, r, redirectURL.Value, http.StatusFound)
	}
}

// CallbackHandler completes Oauth flow and redirects to specified url
func CallbackHandler(auth *Auth) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthFlowSession, err := sessionStore.Get(r, r.FormValue(oauthStateKey))
		handleError(w, err)

		redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
		if !ok {
			handleError(w, errors.New("invalid state parameter"))
		}

		code := r.FormValue(oauthCodeKey)

		tok, err := auth.config.Exchange(context.Background(), code)
		handleError(w, err)

		session, err := sessionStore.New(r, defaultSessionID)
		handleError(w, err)

		session.Values[oauthTokenSessionKey] = tok

		http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
	}
}
