package examples

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/comvex-jp/uipath-go"
	"github.com/comvex-jp/uipath-go/configs"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type Examples struct {
	Client *uipath.Client
}

const folderID = 1234567 // UIPATH Folder ID spcific to a Portal
const username string = "{{UserName}}"
const password string = "{{Password}}"

func Run() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	e := &Examples{
		Client: &uipath.Client{
			HttpClient: &http.Client{Transport: tr},
			Credentials: uipath.Credentials{
				ClientID:          "{{test_client_Id}}", // Deprecated:
				UserKey:           "{{test_user_key}}",  // Deprecated:
				TenantName:        "{{test_tenant_name}}",
				ApplicationID:     "{{test_application_id}}",
				ApplicationSecret: "{{test_application_secret}}",
				Scopes:            "{{test_application_scopes}}",
			},
			BaseURL: "{{test_base_url}}", // UIPATH url specific to the organization/tenant eg. uipath.com/orgName/tenantName/odata
			Cache:   cache.New(5*time.Minute, 10*time.Minute),
		},
	}

	fmt.Println(e.getOauthToken())
	fmt.Println(e.StoreAsset())
	fmt.Println(e.GetAssetById())
	fmt.Println(e.UpdateAsset())
	fmt.Println(e.DeleteAsset())
	fmt.Println(e.ListAssets())
	fmt.Println(e.StoreQueueItem())
	fmt.Println(e.GetQueueItemByID())
	fmt.Println(e.ListQueueItems())
	fmt.Println(e.StoreCredentialVerificationQueueItem())
	fmt.Println(e.StoreDataExtractVerificationQueueItem())
}

func (e *Examples) getOauthToken() uipath.OauthTokenResponse {
	response, err := uipath.GetOAuthToken(e.Client)
	if err != nil {
		panic(err)
	}

	return response
}

func (e *Examples) GetAssetById() (uipath.Asset, error) {
	fmt.Println(e.Client.Cache.Get(configs.UIPathOauthToken))
	asset, err := e.StoreAsset()
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	return aHandler.GetByID(asset.ID)
}

func (e *Examples) ListAssets() ([]uipath.Asset, int, error) {
	var assetList uipath.AssetList

	filters := map[string]string{
		"$top": "1",
	}

	_, err := e.StoreAsset()
	if err != nil {
		return assetList.Value, assetList.Count, err
	}

	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	return aHandler.List(filters)
}

func (e *Examples) UpdateAsset() (uipath.Asset, error) {
	asset, err := e.StoreAsset()
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	updateAsset := uipath.Asset{
		ID:          asset.ID,
		Name:        "NewAssetName",
		ValueType:   "Text",
		StringValue: "Eyyyyyyy",
	}

	return aHandler.Update(updateAsset)
}

func (e *Examples) StoreLoginAsset() (uipath.Asset, error) {
	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	asset := uipath.Asset{
		Name:               "dgm-1-1-credential",
		ValueType:          uipath.ValueTypeCredential,
		ValueScope:         uipath.ValueScopeGlobal,
		CredentialUsername: username,
		CredentialPassword: password,
	}

	return aHandler.Store(asset)
}

func (e *Examples) StoreAsset() (uipath.Asset, error) {
	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	asset := uipath.Asset{
		Name:        fmt.Sprintf("Asset %d", r1.Intn(100)),
		ValueScope:  "Global",
		ValueType:   "Text",
		StringValue: "TestValue",
	}

	return aHandler.Store(asset)
}

func (e *Examples) GetQueueItemByID() (uipath.QueueItem, error) {
	qItem, err := e.StoreQueueItem()
	if err != nil {
		return qItem, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	return queueHandler.GetByID(qItem.ID)
}

func (e *Examples) ListQueueItems() ([]uipath.QueueItem, int, error) {
	var queueItemList uipath.QueueItemList
	filters := map[string]string{
		"$top": "1",
	}

	_, err := e.StoreQueueItem()
	if err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	return queueHandler.List(filters)
}

func (e *Examples) StoreCredentialVerificationQueueItem() (uipath.QueueItem, error) {
	qHandler := uipath.QueueItemHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	// now := time.Now().Format("2006-01-02T15:04:05.4407392Z")
	qI := uipath.QueueItem{
		Priority: uipath.PriorityNormal,
		Name:     "CredentialVerification",
		SpecificContent: map[string]interface{}{
			"CredentialName": "dgm-2-17-credential",
			"AccountID":      2,
			"PortalID":       17,
		},
		Reference: uuid.NewString(),
	}

	return qHandler.Store(qI)
}

func (e *Examples) StoreDataExtractVerificationQueueItem() (uipath.QueueItem, error) {
	qHandler := uipath.QueueItemHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	qI := uipath.QueueItem{
		Name: "DataExtraction",
		SpecificContent: map[string]interface{}{
			"CredentialName":                  "dgm-2-17-credential",
			"AccountID":                       2,
			"PortalID":                        17,
			"SubmittedAt":                     "2022-09-20 10:04:05",
			"LastSubmittedProviderIdentifier": "u1223316",
		},
		Priority:  uipath.PriorityNormal,
		Reference: uuid.NewString(),
	}

	return qHandler.Store(qI)
}

func (e *Examples) StoreQueueItem() (uipath.QueueItem, error) {
	qHandler := uipath.QueueItemHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	// now := time.Now().Format("2006-01-02T15:04:05.4407392Z")
	qI := uipath.QueueItem{
		DueDate:  "2022-04-08T05:37:00.4407392Z",
		Priority: uipath.PriorityNormal,
		Name:     "CredentialVerification",
		SpecificContent: map[string]interface{}{
			"FirstName":   "FirstName Test",
			"LastName":    "LastName Test",
			"Credentials": "cliff-staging-credential",
		},
		Reference: "Petstore",
	}

	return qHandler.Store(qI)
}

func (e *Examples) DeleteAsset() error {
	asset, err := e.StoreAsset()
	if err != nil {
		return err
	}

	aHandler := uipath.AssetHandler{
		Client:   e.Client,
		FolderId: folderID,
	}

	return aHandler.DeleteByID(asset.ID)
}
