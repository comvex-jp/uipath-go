package examples

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/comvex-jp/uipath-go"
	"github.com/comvex-jp/uipath-go/configs"
	"github.com/jarcoal/httpmock"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExamplesTestSuite struct {
	suite.Suite
	c *uipath.Client
}

func (suite *ExamplesTestSuite) SetupTest() {
	cache := cache.New(5*time.Minute, 5*time.Minute)
	suite.c = &uipath.Client{
		HttpClient: &http.Client{Transport: httpmock.DefaultTransport},
		Credentials: uipath.Credentials{
			ClientID:          "testClientID",
			UserKey:           "testUserKey",
			TenantName:        "TestTenant",
			ApplicationID:     "TEST_APP_ID",
			ApplicationSecret: "TEST_APP_SECRET",
			Scopes:            "all",
		},
		BaseURL: "https://cloud.uipath.com/exampleOrg/exampleTenant/odata/",
		Cache:   cache,
	}

	suite.c.Cache.Set(configs.UIPathOauthToken, "=testToken=", 5*time.Minute)
}

func (suite *ExamplesTestSuite) TeardownTest() {
	suite.c.Cache.Flush()
}

func TestExamples(t *testing.T) {
	suite.Run(t, new(ExamplesTestSuite))
}

func (suite *ExamplesTestSuite) TestGetAsset() {
	examples := Examples{
		Client: suite.c,
	}

	asset := map[string]interface{}{
		"Id":          1,
		"Name":        "testGet",
		"ValueScope":  "Global",
		"ValueType":   "Text",
		"StringValue": "TestValue",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.PrepareUIPathResponder(asset, uipath.AssetEndpoint, "POST", 201)
	suite.PrepareUIPathResponder(asset, fmt.Sprintf("%s(%d)", uipath.AssetEndpoint, 1), "GET", 200)

	fetchedAsset, err := examples.GetAssetById()

	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), asset["Name"], fetchedAsset.Name)
	assert.Equal(suite.T(), asset["ValueScope"], fetchedAsset.ValueScope)
	assert.Equal(suite.T(), asset["ValueType"], fetchedAsset.ValueType)
	assert.Equal(suite.T(), asset["StringValue"], fetchedAsset.StringValue)
}

func (suite *ExamplesTestSuite) PrepareUIPathResponder(providedData map[string]interface{}, endpoint string, method string, HTTPStatusCode int) {
	mockResponse, err := json.Marshal(providedData)
	if err != nil {
		panic(err.Error())
	}

	mockURL := suite.c.BaseURL + endpoint

	// Mock responders
	httpmock.RegisterResponder(method, mockURL,
		httpmock.NewStringResponder(HTTPStatusCode, string(mockResponse)))
}
