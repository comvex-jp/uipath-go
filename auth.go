package uipath

import (
	"encoding/json"
)

const OauthURL = "https://account.uipath.com/oauth/token"

// OauthTokenResp is the structure when fetching an oauth token
type OauthTokenResp struct {
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
func GetOAuthToken(c *Client) (OauthTokenResp, error) {
	url := OauthURL

	var result OauthTokenResp

	body := OauthTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.Credentials.ClientID,
		RefreshToken: c.Credentials.UserKey,
	}

	respBody, err := c.Send("POST", url, body, nil, nil)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(respBody, &result); err != nil {
		return result, err
	}

	return result, nil
}
