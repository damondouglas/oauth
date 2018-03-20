package oauth

import (
	"fmt"
	"net/http"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}

// LoginHandler redirects user to oauth flow
func LoginHandler(credentialsPath string, oauthFlowURL string, scopes []string, force bool, offline bool) (func(http.ResponseWriter, *http.Request), error) {
	auth, err := New(credentialsPath, oauthFlowURL, scopes)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := auth.BuildRedirectURL(force, offline)
		sessionID := redirectURL.SessionID
		oauthFlowSession, err := getOauthFlowSession(r, sessionID)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
		}
		redirectURLStr, err := validateRedirectURL(r.FormValue(oauthFlowRedirectKey))
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
		}
		oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURLStr

		if err = oauthFlowSession.Save(r, w); err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
		}
		http.Redirect(w, r, redirectURL.Value, http.StatusFound)
	}, nil
}
