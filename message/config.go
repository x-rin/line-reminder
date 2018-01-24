package message

type lineConfig struct {
	AccessToken	string 	`json:"access_token"`
	ExpiresIn	int		`json:"expires_in"`
	TokenType	string	`json:"token_type"`
}
