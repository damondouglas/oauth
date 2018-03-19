package oauth

import (
	"testing"

	"github.com/damondouglas/assert"
)

func TestConfigureOauthFromData(t *testing.T) {
	a := assert.New(t)
	data := []byte(`{"web": {"client_id": "foo", "client_secret": "bar"}}`)
	config := configureOauthFromData(data)
	a.NotEquals(config, nil)
	a.Equals(config.ClientID, "foo")
	a.Equals(config.ClientSecret, "bar")
}

func TestConfigureOauthFromFilePath(t *testing.T) {
	a := assert.New(t)
	config, err := configureOauthFromFilePath("./mock/client_secret.json")
	a.Equals(err, nil)
	a.NotEquals(config, nil)
	a.Equals(config.ClientID, "some_client_id.apps.googleusercontent.com")
	a.Equals(config.ClientSecret, "some_client_secret")
}
