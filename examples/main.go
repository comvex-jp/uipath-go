package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/comvex-jp/uipath-go"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	credentials := uipath.Credentials{
		ClientID:   "test_key_here",
		UserKey:    "test_user_key_here",
		TenantName: "tenant_name_here",
	}

	httpClient := &http.Client{Transport: tr}

	c := uipath.Client{
		HttpClient:  httpClient,
		Credentials: credentials,
		URLEndpoint: "https://cloud.uipath.com/{org_name}/{tenant_name}/odata/",
	}

	res, err := uipath.GetOAuthToken(&c)
	if err != nil {
		fmt.Println(res, err)
	}

	c.Credentials.Token = res.AccessToken

	fmt.Println(StoreAsset(c))
	fmt.Println(GetAssetById(c))
	fmt.Println(UpdateAsset(c))
	fmt.Println(ListAssets(c))
	fmt.Println(StoreQueueItem(c))
	fmt.Println(GetQueueItemByID(c))
	fmt.Println(ListQueueItems(c))
}

func GetAssetById(c uipath.Client) (uipath.Asset, error) {
	folderID := uint(292388)

	asset, err := StoreAsset(c)
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client: &c,
	}

	return aHandler.GetByID(asset.ID, folderID)
}

func ListAssets(c uipath.Client) ([]uipath.Asset, int, error) {
	var assetList uipath.AssetList

	folderID := uint(292388)

	filters := map[string]string{
		"$top": "1",
	}

	_, err := StoreAsset(c)
	if err != nil {
		return assetList.Value, assetList.Count, err
	}

	aHandler := uipath.AssetHandler{
		Client: &c,
	}

	return aHandler.List(filters, folderID)
}

func UpdateAsset(c uipath.Client) (uipath.Asset, error) {
	folderID := uint(292388)

	asset, err := StoreAsset(c)
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client: &c,
	}

	updateAsset := uipath.Asset{
		ID:          asset.ID,
		Name:        "NewAssetName",
		ValueType:   "Text",
		StringValue: "Eyyyyyyy",
	}

	return aHandler.Update(updateAsset, folderID)
}

func StoreAsset(c uipath.Client) (uipath.Asset, error) {
	folderID := uint(292388)
	aHandler := uipath.AssetHandler{
		Client: &c,
	}

	asset := uipath.Asset{
		Name:        fmt.Sprintf("Asset %d", rand.Intn(100)),
		ValueScope:  "Global",
		ValueType:   "Text",
		StringValue: "TestValue",
	}

	return aHandler.Store(asset, folderID)
}

func GetQueueItemByID(c uipath.Client) (uipath.QueueItem, error) {
	folderID := uint(292388)

	qItem, err := StoreQueueItem(c)
	if err != nil {
		return qItem, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client: &c,
	}

	return queueHandler.GetByID(qItem.ID, folderID)
}

func ListQueueItems(c uipath.Client) ([]uipath.QueueItem, int, error) {
	var queueItemList uipath.QueueItemList
	folderID := uint(292388)

	filters := map[string]string{
		"$top": "1",
	}

	_, err := StoreQueueItem(c)
	if err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client: &c,
	}

	return queueHandler.List(filters, folderID)
}

func StoreQueueItem(c uipath.Client) (uipath.QueueItem, error) {
	folderID := uint(292388)
	qHandler := uipath.QueueItemHandler{
		Client: &c,
	}

	now := time.Now().Format("2006-01-02T15:04:05.4407392Z")
	qI := uipath.QueueItem{
		DeferDate: now,
		DueDate:   now,
		Priority:  uipath.PriorityNormal,
		Name:      "ContactCreation",
		SpecificContent: map[string]interface{}{
			"Test":     "Test from API",
			"TestBool": false,
		},
		Reference: "Petstore",
	}

	return qHandler.Store(qI, folderID)
}
