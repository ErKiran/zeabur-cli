package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	KeyTokenString = "token"
	KeyUser        = "user"

	// if we use "token.access_token" as key, the env ZEABUR_TOKEN will override all token details,
	// because "token" is the prefix of "token.access_token".
	// Therefore, we use "token_detail" as the key to store all token details.

	KeyTokenDetail  = "token_detail"
	KeyTokenAccess  = KeyTokenDetail + ".access_token"
	KeyTokenExpiry  = KeyTokenDetail + ".expiry"
	KeyTokenType    = KeyTokenDetail + ".token_type"
	KeyTokenRefresh = KeyTokenDetail + ".refresh_token"
)

const (
	KeyInteractive      = "interactive"
	KeyAutoRefreshToken = "auto_refresh_token"
)

type Config interface {
	GetTokenString() string // token string is the single token string, it may be set by user or generated from OAuth2
	SetTokenString(token string)

	GetToken() *oauth2.Token // token is the detail of token, if the token is from OAuth2, it will be set
	SetToken(token *oauth2.Token)

	GetUser() string
	SetUser(user string)

	Write() error
}

type config struct {
	path string
}

func New(path string) Config {
	// create the config file and init viper
	initViper(path)

	return &config{
		path: path,
	}
}

func (c *config) GetTokenString() string {
	return viper.GetString(KeyTokenString)
}

func (c *config) SetTokenString(token string) {
	viper.Set(KeyTokenString, token)
}

func (c *config) GetToken() *oauth2.Token {
	token := &oauth2.Token{}
	token.AccessToken = viper.GetString(KeyTokenAccess)
	token.RefreshToken = viper.GetString(KeyTokenRefresh)
	token.TokenType = viper.GetString(KeyTokenType)
	token.Expiry = viper.GetTime(KeyTokenExpiry)
	if token.AccessToken == "" || token.RefreshToken == "" || token.TokenType == "" || token.Expiry.IsZero() {
		return nil
	}
	return token
}

func (c *config) SetToken(token *oauth2.Token) {
	if token == nil {
		viper.Set(KeyTokenDetail, "")
		return
	}
	viper.Set(KeyTokenAccess, token.AccessToken)
	viper.Set(KeyTokenRefresh, token.RefreshToken)
	viper.Set(KeyTokenType, token.TokenType)
	viper.Set(KeyTokenExpiry, token.Expiry)
}

func (c *config) GetUser() string {
	return viper.GetString(KeyUser)
}

func (c *config) SetUser(user string) {
	viper.Set(KeyUser, user)
}

func (c *config) Write() error {
	return viper.WriteConfig()
}

var _ Config = &config{}