package uipath

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	c *Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.c = &Client{
		HttpClient: &httpClientMock{},
	}
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type httpClientMock struct {
	MockedDo func(req *http.Request) (*http.Response, error)
}

func (c *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.MockedDo(req)
}

func (suite *ClientTestSuite) TestGetOathToken() {
	suite.c = &Client{
		HttpClient: &http.Client{Transport: httpmock.DefaultTransport},
		Credentials: Credentials{
			ClientID: "asdasdasd",
			UserKey:  "asdasdasd",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.PrepareUIPathAuthAPIResponder(PrepareOauthTokenData(), OauthURL, "", "POST", 201, true)

	resp, _ := GetOAuthToken(suite.c)

	assert.Equal(suite.T(), "ey0123456789", resp.AccessToken)
}

func (suite *ClientTestSuite) TestSend() {
	type bodyMock struct {
		Message string `json:"message"`
	}

	resBody := ``

	body := ioutil.NopCloser(bytes.NewReader([]byte(resBody)))

	suite.c.HttpClient = &httpClientMock{
		MockedDo: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       body,
				},
				nil
		},
	}

	res, err := suite.c.Send("GET", "/hoge", bodyMock{Message: "test"}, nil, nil)

	var test map[string]interface{}

	json.Unmarshal(res, &test)
	suite.T().Log(test, err)

	assert.Equal(suite.T(), "HTTP Error 401: Unauthorized", err.Error())
	assert.Equal(suite.T(), string(res), resBody)
}

func PrepareOauthTokenData() OauthTokenResp {
	return OauthTokenResp{
		TokenType:   "Bearer",
		IDToken:     "ey0123456789",
		ExpiresIn:   2678400,
		AccessToken: "ey0123456789",
		Scope:       "all",
	}
}

func (suite *ClientTestSuite) TestSendWithAuthorization() {
	resBody := `{"http_status_code": 200}`

	suite.c.Credentials.Token = "abcd1234"

	suite.c.HttpClient = &httpClientMock{
		MockedDo: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(resBody))),
				},
				nil
		},
	}

	res, err := suite.c.SendWithAuthorization("GET", "/", nil, nil, nil)

	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), string(res), resBody)
}

func (suite *ClientTestSuite) PrepareUIPathAuthAPIResponder(providedData OauthTokenResp, envBaseUrl string, endpoint string, method string, HTTPStatusCode int, success bool) {

	mockResponse, err := json.Marshal(providedData)
	if err != nil {
		panic(err.Error())
	}

	mockURL := envBaseUrl + endpoint

	// Mock responders
	httpmock.RegisterResponder(method, mockURL,
		httpmock.NewStringResponder(HTTPStatusCode, string(mockResponse)))
}
