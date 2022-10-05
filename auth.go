package uipath

import (
	"encoding/json"
	"errors"
)

// Deprecated: User key based authentication is deprecated since October, 2021
// See more details in https://www.uipath.com/ja/resources/knowledge-base/implementing-orchestrator-api-with-oauth
const DeprecatedOauthURL = "https://account.uipath.com/oauth/token"

const OauthURL = "https://cloud.uipath.com/identity_/connect/token"

// OauthTokenResponse is the structure when fetching an oauth token
type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Deprecated: OauthTokenRequest are the basic values needed to get an oauth token
type DeprecatedOauthTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

// OauthTokenRequest are the basic values needed to get an oauth token
type OauthTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

// Deprecated: getOAuthToken requests oauthtoken using UserKey
func DeprecatedGetOAuthToken(c *Client) (OauthTokenResponse, error) {
	url := DeprecatedOauthURL

	var result OauthTokenResponse

	body := DeprecatedOauthTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.Credentials.ClientID,
		RefreshToken: c.Credentials.UserKey,
	}

	respBody, err := c.Send("POST", url, body, map[string]string{}, nil)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(respBody, &result)

	return result, err
}

// GetOAuthToken requests oauthtoken using ClientCredentials
func GetOAuthToken(c *Client) (OauthTokenResponse, error) {
	var result OauthTokenResponse

	body := OauthTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.Credentials.ApplicationID,
		ClientSecret: c.Credentials.ApplicationSecret,
		Scope:        c.Credentials.Scopes,
	}

	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	respBody, err := c.Send("POST", OauthURL, body, headers, nil)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(respBody, &result)

	if result.AccessToken == "" {
		return result, errors.New("Empty Access Token Error")
	}

	return result, err
}
