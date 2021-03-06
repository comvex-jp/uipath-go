package uipath

import (
	"encoding/json"
)

const OauthURL = "https://account.uipath.com/oauth/token"

// OauthTokenResponse is the structure when fetching an oauth token
type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// OauthTokenRequest are the basic values needed to get an oauth token
type OauthTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

// GetOAuthToken requests oauthtoken using credentials
func GetOAuthToken(c *Client) (OauthTokenResponse, error) {
	url := OauthURL

	var result OauthTokenResponse

	body := OauthTokenRequest{
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
