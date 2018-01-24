package message

type lineConfig struct {
	AccessToken	string 	`json:"access_token"`
	ExpiresIn	int		`json:"expires_in"`
	TokenType	string	`json:"token_type"`
}

type LineConfig interface {
	GetLineConfig(accessToken string, expireIn int, tokenType string) *lineConfig
}

func GetLineConfig(accessToken string, expireIn int, tokenType string) *lineConfig {
	return &lineConfig{
		AccessToken: accessToken,
		ExpiresIn: expireIn,
		TokenType: tokenType,
	}
}
