package api

import (
	"log"
	"net/url"
	"os"
	"net/http"
	"strings"
	"encoding/json"
)

type LineConfig struct {
	AccessToken	string 	`json:"access_token"`
	ExpiresIn	int		`json:"expires_in"`
	TokenType	string	`json:"token_type"`
}

func NewLineConfig() *LineConfig{
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", os.Getenv("CHANNEL_ID"))
	values.Set("client_secret", os.Getenv("CHANNEL_SECRET"))

	url := "https://api.line.me/v2/oauth/accessToken"

	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	config := new(LineConfig)

	if err := json.NewDecoder(res.Body).Decode(config); err != nil {
		log.Fatal(err)
	}

	return config
}
