package uipath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/comvex-jp/uipath-go/configs"
	"github.com/patrickmn/go-cache"
)

const (
	HeaderOrganizationUnitId = "X-UIPATH-OrganizationUnitId"
	HeaderAuthorization      = "Authorization"
	HeaderTenantName         = "X-UIPATH-TenantName"
)

var httpSuccessCodes = map[int]string{
	200: "success",
	201: "created",
	204: "no content",
}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines how the UIPath client looks like
type Client struct {
	HttpClient  HttpClientInterface
	Credentials Credentials
	BaseURL     string
	Cache       *cache.Cache
}

// Credentials struct defines what items are needed for the client credentials
type Credentials struct {
	ClientID          string // Deprecated: Use ApplicationID
	UserKey           string // Deprecated: Use ApplicationSecret
	TenantName        string
	Token             string
	ApplicationID     string
	ApplicationSecret string
	Scopes            string
}

type resultCode struct {
	Result string `json:"result"`
}

// GetAuthHeaderValue gets the token if it exists and fetches if it does not
func (client *Client) GetAuthHeaderValue() (string, error) {
	var token string

	res, found := client.Cache.Get(configs.UIPathOauthToken)
	if !found {
		fetchedTokenData, err := GetOAuthToken(client)
		if err != nil {
			log.Println("Error fetching access token: ", err.Error())

			fetchedTokenData, err = DeprecatedGetOAuthToken(client)
		}

		if err != nil {
			return token, err
		}

		token = fetchedTokenData.AccessToken
		expiresIn := fetchedTokenData.ExpiresIn

		expireInStr := strconv.Itoa(expiresIn) + "s"

		parsedExpiresIn, err := time.ParseDuration(expireInStr)
		if err != nil {
			return token, err
		}

		client.Cache.Set(configs.UIPathOauthToken, token, parsedExpiresIn)
		return token, nil
	}

	return res.(string), nil
}

// Send handles all requests going out for uipath clinet
func (client Client) Send(requestMethod string, url string, body interface{}, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return jsonBody, err
	}

	req, err := http.NewRequest(requestMethod, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return jsonBody, err
	}

	if len(queryParams) > 0 {
		attachQueryParams(req, queryParams)
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	headers[HeaderTenantName] = client.Credentials.TenantName

	attachHeaders(req, headers)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return jsonBody, err
	}

	fmt.Printf("response: %+v\n", resp)
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("body: %+v\n", string(bodyResp))

	defer func(Body io.ReadCloser, Request *http.Response) {
		if err := Body.Close(); err != nil {
			log.Println("Error closing body from response: ", resp)
		}
	}(resp.Body, resp)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jsonBody, err
	}

	// Handle any errors from the response
	if _, ok := httpSuccessCodes[resp.StatusCode]; !ok {
		return respBody, ErrorResponseHandler(resp.StatusCode, respBody)
	}

	return respBody, err
}

// SendWithAuthorization attaches the authorization token to the headers and then completes the request
func (client *Client) SendWithAuthorization(requestMethod, url string, body interface{}, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return jsonBody, err
	}

	token, err := client.GetAuthHeaderValue()
	if err != nil {
		return jsonBody, err
	}

	headers[HeaderAuthorization] = "Bearer " + token

	return client.Send(requestMethod, url, body, headers, queryParams)
}

func attachHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func attachQueryParams(req *http.Request, queryParams map[string]string) {
	q := req.URL.Query()

	for i, v := range queryParams {
		q.Add(i, v)
	}

	req.URL.RawQuery = q.Encode()
}
