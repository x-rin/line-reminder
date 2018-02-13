package reminder

import (
	"github.com/gin-gonic/gin/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// AccessConfig - Line APIにアクセスするためのConfig
type AccessConfig struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

// ChannelConfig - Line Channelに関するConfig
type ChannelConfig struct {
	ID     string
	Secret string
}

// Config - AccessConfig, ChannelConfigを統合したもの
type Config struct {
	Access  *AccessConfig
	Channel *ChannelConfig
}

// NewConfig - Configを生成
func NewConfig() (*Config, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", os.Getenv("CHANNEL_ID"))
	values.Set("client_secret", os.Getenv("CHANNEL_SECRET"))

	requestURL := "https://api.line.me/v2/oauth/accessToken"
	req, err := http.NewRequest(
		"POST",
		requestURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var accessConfig *AccessConfig
	if err := json.NewDecoder(res.Body).Decode(accessConfig); err != nil {
		return nil, err
	}
	return &Config{
		Access: accessConfig,
		Channel: &ChannelConfig{
			ID:     os.Getenv("CHANNEL_ID"),
			Secret: os.Getenv("CHANNEL_SECRET"),
		},
	}, nil
}
