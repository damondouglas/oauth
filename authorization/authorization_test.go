package authorization

import (
	"net/http/httptest"
	"testing"

	"github.com/damondouglas/assert"
	"github.com/damondouglas/oauth"
)

func TestHandle(t *testing.T) {
	a := assert.New(t)
	redirectURL := "https://oauth-redirect.googleusercontent.com/r/my_project_id"
	scopes := []string{"email", "profile"}
	url := "https://myservice.example.com/auth?client_id=some_client_id.apps.googleusercontent.com&redirect_uri=https://oauth-redirect.googleusercontent.com/r/my_project_id&state=STATE_STRING&scope=email%20profile&response_type=code"
	r := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	oauthConfig, err := oauth.ConfigFromPath("../mock/client_secret.json", redirectURL, scopes)
	a.HandleError(err)
	auth := New(oauthConfig, redirectURL, scopes, true, true)
	auth.Handle(w, r)
	resp := w.Result()
	location, err := resp.Location()
	a.HandleError(err)
	a.Equals(location.String(), "https://accounts.google.com/o/oauth2/auth?access_type=offline&approval_prompt=force&client_id=some_client_id.apps.googleusercontent.com&redirect_uri=https%3A%2F%2Foauth-redirect.googleusercontent.com%2Fr%2Fmy_project_id&response_type=code&scope=email+profile&state=STATE_STRING")
}
