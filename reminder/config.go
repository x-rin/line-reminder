package reminder

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type authResponse struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

// GetToken - ChannelTokenを取得する
func GetChannelToken(clientID, clientSecret string) (*string, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", clientID)
	values.Set("client_secret", clientSecret)
	requestURL := "https://api.line.me/v2/oauth/accessToken"
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var authResponse authResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&authResponse); err != nil {
		return nil, err
	}

	return &authResponse.Token, nil
}
