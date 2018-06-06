package query

import (
	"sort"
	"testing"

	"github.com/damondouglas/assert"
)

func TestParseAuthorizationQuery(t *testing.T) {
	a := assert.New(t)
	queryData := "client_id=GOOGLE_CLIENT_ID&redirect_uri=REDIRECT_URI&state=STATE_STRING&scope=email%20profile&response_type=code"
	qry, err := ParseAuthorizationQuery(queryData)
	a.HandleError(err)
	a.Equals(qry.IsValid(), true)
	a.Equals(qry.ClientID(), "GOOGLE_CLIENT_ID")
	a.Equals(qry.RedirectURI(), "REDIRECT_URI")
	a.Equals(qry.State(), "STATE_STRING")
	sort.Strings(qry.scope)
	a.Equals(len(qry.Scope()), 2)
	scope := []string{"email", "profile"}
	for i, s := range scope {
		a.Equals(qry.Scope()[i], s)
	}
	a.Equals(qry.ResponseType(), "code")
}

func TestParseTokenQuery(t *testing.T) {
	a := assert.New(t)
	authorizationQueryData := "client_id=GOOGLE_CLIENT_ID&client_secret=GOOGLE_CLIENT_SECRET&grant_type=authorization_code&code=AUTHORIZATION_CODE"
	qry, err := ParseTokenQuery(authorizationQueryData)
	a.HandleError(err)
	a.Equals(qry.IsValid(), true)
	a.Equals(qry.ClientID(), "GOOGLE_CLIENT_ID")
	a.Equals(qry.ClientSecret(), "GOOGLE_CLIENT_SECRET")
	a.Equals(qry.GrantType(), "authorization_code")
	a.Equals(qry.Code(), "AUTHORIZATION_CODE")
	a.Equals(qry.Mode, AuthorizationCodeMode)

	refreshQueryData := "client_id=GOOGLE_CLIENT_ID&client_secret=GOOGLE_CLIENT_SECRET&grant_type=refresh_token&refresh_token=REFRESH_TOKEN"
	qry, err = ParseTokenQuery(refreshQueryData)
	a.HandleError(err)
	a.Equals(qry.IsValid(), true)
	a.Equals(qry.ClientID(), "GOOGLE_CLIENT_ID")
	a.Equals(qry.ClientSecret(), "GOOGLE_CLIENT_SECRET")
	a.Equals(qry.GrantType(), "refresh_token")
	a.Equals(qry.RefreshToken(), "REFRESH_TOKEN")
	a.Equals(qry.Mode, RefreshTokenMode)
}
