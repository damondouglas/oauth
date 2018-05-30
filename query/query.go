package query

import (
	"net/url"
	"strings"
)

const (
	stateKey                 = "state"
	codeKey                  = "code"
	scopeKey                 = "scope"
	redirectURIKey           = "redirect_uri"
	responseTypeKey          = "response_type"
	authorizationCodeModeKey = 0
	authorizationCodeKey     = "authorization_code"
	refreshTokenModeKey      = 1
	refreshTokenKey          = "refresh_token"
	clientIDKey              = "client_id"
	clientSecretKey          = "client_secret"
	grantTypeKey             = "grant_type"
)

var (
	// AuthorizationCodeMode specifies grant_type=authorization_code
	AuthorizationCodeMode *TokenMode

	// RefreshTokenMode specifies grant_type=refresh_token
	RefreshTokenMode *TokenMode
)

func init() {
	AuthorizationCodeMode = new(TokenMode)
	AuthorizationCodeMode.key = authorizationCodeModeKey

	RefreshTokenMode = new(TokenMode)
	RefreshTokenMode.key = refreshTokenModeKey
}

type query struct {
	clientID     string
	clientSecret string
	code         string
	grantType    string
	redirectURI  string
	refreshToken string
	responseType string
	scope        []string
	state        string
}

// AuthorizationQuery represents parsed query string from OAuth2 authorization.
type AuthorizationQuery struct {
	query
}

// ParseAuthorizationQuery parses ?client_id=GOOGLE_CLIENT_ID&redirect_uri=REDIRECT_URI&state=STATE_STRING&scope=REQUESTED_SCOPES&response_type=code.
// It returns AuthorizationQuery.
func ParseAuthorizationQuery(queryData string) (*AuthorizationQuery, error) {
	var s *AuthorizationQuery
	q, err := parse(queryData)
	if err != nil {
		return nil, err
	}

	temp := AuthorizationQuery{q}
	s = &temp

	return s, err
}

// IsValid returns true if the query string is properly parsed.
func (s *AuthorizationQuery) IsValid() bool {
	hasData := allNotEmpty(
		[]string{
			s.clientID,
			s.redirectURI,
			s.state,
		},
	)
	hasData = hasData && len(s.scope) > 0

	return hasData && s.responseType == codeKey
}

// RedirectURI returns value parsed from redirect_uri=REDIRECT_URI
func (s *AuthorizationQuery) RedirectURI() string {
	return s.redirectURI
}

// State returns value parsed from state=STATE_STRING
func (s *AuthorizationQuery) State() string {
	return s.state
}

// Scope returns value parsed from scope=REQUESTED_SCOPES
func (s *AuthorizationQuery) Scope() []string {
	return s.scope
}

// ResponseType returns value parsed from response_type=code
func (s *AuthorizationQuery) ResponseType() string {
	return s.responseType
}

// TokenQuery represents parsed form data from OAuth2 token exchange requests.
type TokenQuery struct {
	query
	Mode *TokenMode
}

// ParseTokenQuery parses form data from OAuth2 token exchange requests.
// Query string is of the form client_id=GOOGLE_CLIENT_ID&client_secret=GOOGLE_CLIENT_SECRET&<token_request_detail>.
// The part <token_request_detail> may be:
// &grant_type=authorization_code&code=AUTHORIZATION_CODE or
// &grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func ParseTokenQuery(queryData string) (*TokenQuery, error) {
	var t *TokenQuery
	q, err := parse(queryData)
	if err != nil {
		return nil, err
	}
	var mode *TokenMode
	switch q.grantType {
	case authorizationCodeKey:
		mode = AuthorizationCodeMode
	case refreshTokenKey:
		mode = RefreshTokenMode
	}
	temp := TokenQuery{q, mode}
	t = &temp

	return t, err
}

// IsValid returns true if query string was properly parsed.
func (t *TokenQuery) IsValid() bool {
	hasData := allNotEmpty(
		[]string{
			t.clientID,
			t.clientSecret,
			t.grantType,
		},
	)
	hasData = hasData && t.Mode != nil
	hasData = hasData && (t.refreshToken != "" || t.code != "")

	return hasData
}

// ClientSecret returns value parsed from client_secret=GOOGLE_CLIENT_SECRET
func (t *TokenQuery) ClientSecret() string {
	return t.clientSecret
}

// Code returns value parsed from code=AUTHORIZATION_CODE
func (t *TokenQuery) Code() string {
	return t.code
}

// RefreshToken returns value parsed from refresh_token=REFRESH_TOKEN
func (t *TokenQuery) RefreshToken() string {
	return t.refreshToken
}

// GrantType returns value parsed from grant_type=(authorization_code|refresh_token)
func (t *TokenQuery) GrantType() string {
	return t.grantType
}

// ClientID returns value parsed from client_id=CLIENT_ID
func (q *query) ClientID() string {
	return q.clientID
}

func (q *query) setClientID(qry url.Values) {
	value := qry.Get(clientIDKey)
	if value != "" {
		q.clientID = value
	}
}

func (q *query) setCode(qry url.Values) {
	value := qry.Get(codeKey)
	if value != "" {
		q.code = value
	}
}

func (q *query) setGrantType(qry url.Values) {
	value := qry.Get(grantTypeKey)
	if value != "" {
		q.grantType = value
	}
}

func (q *query) setRedirectURI(qry url.Values) {
	value := qry.Get(redirectURIKey)
	if value != "" {
		q.redirectURI = value
	}
}

func (q *query) setRefreshToken(qry url.Values) {
	value := qry.Get(refreshTokenKey)
	if value != "" {
		q.refreshToken = value
	}
}

func (q *query) setResponseType(qry url.Values) {
	value := qry.Get(responseTypeKey)
	if value != "" {
		q.responseType = value
	}
}

func (q *query) setClientSecret(qry url.Values) {
	value := qry.Get(clientSecretKey)
	if value != "" {
		q.clientSecret = value
	}
}

func (q *query) setState(qry url.Values) {
	value := qry.Get(stateKey)
	if value != "" {
		q.state = value
	}
}

func (q *query) setScope(qry url.Values) {
	value := qry.Get(scopeKey)
	if value != "" {
		q.scope = strings.Split(value, " ")
	}
}

func parse(queryData string) (query, error) {
	var q query
	qry, err := url.ParseQuery(queryData)
	if err != nil {
		return q, err
	}

	q.setClientID(qry)
	q.setClientSecret(qry)
	q.setCode(qry)
	q.setGrantType(qry)
	q.setRedirectURI(qry)
	q.setRefreshToken(qry)
	q.setResponseType(qry)
	q.setScope(qry)
	q.setState(qry)

	return q, nil
}

// TokenMode specifies grant_type=authorization_code or refresh_token
type TokenMode struct {
	key int
}

// String returns authorization_code or refresh_token
func (mode *TokenMode) String() string {
	keys := [...]string{
		authorizationCodeKey,
		refreshTokenKey,
	}
	return keys[mode.key]
}

func allNotEmpty(tokens []string) bool {
	for _, t := range tokens {
		if t == "" {
			return false
		}
	}
	return true
}
