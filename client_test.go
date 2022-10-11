package uipath

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/comvex-jp/uipath-go/configs"
	"github.com/jarcoal/httpmock"
	"github.com/patrickmn/go-cache"
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

	suite.c.Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func (suite *ClientTestSuite) TearDownTest() {
	suite.c.Cache.Flush()
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func PrepareOauthTokenData() OauthTokenResponse {
	return OauthTokenResponse{
		TokenType:   "Bearer",
		IDToken:     "ey0123456789",
		ExpiresIn:   3600,
		AccessToken: "ey0123456789",
		Scope:       "all",
	}
}

type httpClientMock struct {
	MockedDo func(req *http.Request) (*http.Response, error)
}

func (c *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.MockedDo(req)
}

func (suite *ClientTestSuite) TestGetOathToken() {
	suite.c.HttpClient = &http.Client{Transport: httpmock.DefaultTransport}
	suite.c.Credentials = Credentials{
		ClientID:          "asdasdasd",
		UserKey:           "asdasdasd",
		ApplicationID:     "TEST_APP_ID",
		ApplicationSecret: "TEST_APP_SECRET",
		Scopes:            "all",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.PrepareUIPathAuthAPIResponder(PrepareOauthTokenData(), OauthURL, "", "POST", 201)

	resp, _ := GetOAuthToken(suite.c)

	assert.Equal(suite.T(), "ey0123456789", resp.AccessToken)
}

func (suite *ClientTestSuite) TestGetOathTokenFallback() {
	suite.c.HttpClient = &http.Client{Transport: httpmock.DefaultTransport}
	suite.c.Credentials = Credentials{
		ClientID:          "asdasdasd",
		UserKey:           "asdasdasd",
		ApplicationID:     "asdasdasd",
		ApplicationSecret: "asdasdasd",
		Scopes:            "all",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.PrepareUIPathAuthAPIResponder(OauthTokenResponse{}, OauthURL, "", "POST", 404)
	suite.PrepareUIPathAuthAPIResponder(PrepareOauthTokenData(), DeprecatedOauthURL, "", "POST", 201)

	token, err := suite.c.GetAuthHeaderValue()
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "ey0123456789", token)
}

func (suite *ClientTestSuite) TestGetCachedAuthHeaderValue() {
	testAuthToken := "=testToken="
	suite.c.Cache.Set(configs.UIPathOauthToken, testAuthToken, 1*time.Minute)

	token, _ := suite.c.GetAuthHeaderValue()

	assert.Equal(suite.T(), testAuthToken, token)
}

func (suite *ClientTestSuite) TestGetCachedAuthHeaderValueWhenTokenIsEmpty() {
	suite.c.HttpClient = &http.Client{Transport: httpmock.DefaultTransport}
	suite.c.Credentials = Credentials{
		ClientID: "asdasdasd",
		UserKey:  "asdasdasd",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.PrepareUIPathAuthAPIResponder(PrepareOauthTokenData(), OauthURL, "", "POST", 201)

	token, _ := suite.c.GetAuthHeaderValue()

	assert.Equal(suite.T(), "ey0123456789", token)
}

func (suite *ClientTestSuite) TestSend() {
	header := map[string]string{}
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

	res, err := suite.c.Send("GET", "/hoge", bodyMock{Message: "test"}, header, nil)

	var test map[string]interface{}

	json.Unmarshal(res, &test)

	assert.Equal(suite.T(), "HTTP Error 401: Unauthorized", err.Error())
	assert.Equal(suite.T(), string(res), resBody)
}

func (suite *ClientTestSuite) TestSendWithAuthorization() {
	header := map[string]string{}
	resBody := `{"http_status_code": 200}`

	suite.c.HttpClient = &httpClientMock{
		MockedDo: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(resBody))),
				},
				nil
		},
	}

	suite.c.Cache.Set(configs.UIPathOauthToken, "=testToken=", 1*time.Minute)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	res, err := suite.c.SendWithAuthorization("GET", "/", nil, header, nil)

	assert.Equal(suite.T(), err, nil)
	assert.Equal(suite.T(), string(res), resBody)
}

func (suite *ClientTestSuite) PrepareUIPathAuthAPIResponder(providedData OauthTokenResponse, envBaseUrl string, endpoint string, method string, HTTPStatusCode int) {

	mockResponse, err := json.Marshal(providedData)
	if err != nil {
		panic(err.Error())
	}

	mockURL := envBaseUrl + endpoint

	// Mock responders
	httpmock.RegisterResponder(method, mockURL,
		httpmock.NewStringResponder(HTTPStatusCode, string(mockResponse)))
}
