package uipath

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const OrchestratorURL = "https:://orchestrator-url.com/"

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines how the UIPath client looks like
type Client struct {
	HttpClient                 HttpClientInterface
	Credentials                Credentials
	URLEndpoint                string
	Prefix                     string
	FailedUnauthorizedAttempts uint8
}

// Credentials struct defines what items are needed for the client credentials
type Credentials struct {
	ClientID   string
	UserKey    string
	TenantName string
	Token      string
}

type resultCode struct {
	Result string `json:"result"`
}

// Send handles all requests going out for uipath clinet
func (c Client) Send(requestMethod string, url string, body interface{}, headers map[string]string, queryParams map[string]string) ([]byte, error) {
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-UIPATH-TenantName", c.Credentials.TenantName)

	attachHeaders(req, headers)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return jsonBody, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	// Handle any errors from the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return respBody, ErrorResponseHandler(resp.StatusCode, respBody)
	}

	return respBody, err
}

// SendWithAuthorization attaches the authorization token to the headers and then completes the request
func (c *Client) SendWithAuthorization(requestMethod, url string, body interface{}, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	headersWithAuthorization := make(map[string]string)

	for k, v := range headers {
		headersWithAuthorization[k] = v
	}

	headersWithAuthorization["Authorization"] = "Bearer " + c.Credentials.Token

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return jsonBody, err
	}

	return c.Send(requestMethod, url, body, headersWithAuthorization, queryParams)
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
