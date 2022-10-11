package uipath

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// GetOAuthToken requests Oauth Token using ClientCredentials
func GetOAuthToken(c *Client) (OauthTokenResponse, error) {
	var result OauthTokenResponse

	queryParams := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.Credentials.ApplicationID,
		"client_secret": c.Credentials.ApplicationSecret,
		"scope":         c.Credentials.Scopes,
	}

	form := url.Values{}
	for k, v := range queryParams {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", OauthURL, strings.NewReader(form.Encode()))
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return result, err
	}

	defer func(Body io.ReadCloser, Request *http.Response) {
		if err := Body.Close(); err != nil {
			log.Println("Error closing body from response: ", resp)
		}
	}(resp.Body, resp)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// Handle any errors from the response
	if _, ok := httpSuccessCodes[resp.StatusCode]; !ok {
		return result, ErrorResponseHandler(resp.StatusCode, respBody)
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return result, err
	}

	if result.AccessToken == "" {
		return result, errors.New("Empty Access Token Error")
	}

	return result, nil
}
