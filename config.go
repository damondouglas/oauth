package oauth

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Web *webOauthConfig `json:"web"`
}

type webOauthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func configureOauthFromData(data []byte) *webOauthConfig {
	var c *config
	json.Unmarshal(data, &c)
	return c.Web
}

func configureOauthFromFilePath(pathToCredentials string) (*webOauthConfig, error) {
	data, err := ioutil.ReadFile(pathToCredentials)
	if err != nil {
		return nil, err
	}
	return configureOauthFromData(data), nil
}
