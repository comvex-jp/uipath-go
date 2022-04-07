package examples

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/comvex-jp/uipath-go"
)

func Run() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	credentials := uipath.Credentials{
		ClientID:   "{{ClientID}}",
		UserKey:    "{{UserKey}}",
		TenantName: "{{TenantName}}",
	}

	httpClient := &http.Client{Transport: tr}

	c := uipath.Client{
		HttpClient:  httpClient,
		Credentials: credentials,
		URLEndpoint: "{{URLEndpoint}}",
	}

	res, err := uipath.GetOAuthToken(&c)
	if err != nil {
		fmt.Println(res, err)
	}

	c.Credentials.Token = res.AccessToken

	fmt.Print("Store Asset ")
	fmt.Println(StoreAsset(c))
	fmt.Print("Get Asset ")
	fmt.Println(GetAssetById(c))
	fmt.Print("Update Asset ")
	fmt.Println(UpdateAsset(c))
	fmt.Print("List Assets ")
	fmt.Println(ListAssets(c))
	fmt.Print("Store QueueItem ")
	fmt.Println(StoreQueueItem(c))
	fmt.Print("Get QueueItem ")
	fmt.Println(GetQueueItemByID(c))
	fmt.Print("List QueueItems ")
	fmt.Println(ListQueueItems(c))
}

func GetAssetById(c uipath.Client) (uipath.Asset, error) {
	asset, err := StoreAsset(c)
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client:   &c,
		FolderId: uint(292388),
	}

	return aHandler.GetByID(asset.ID)
}

func ListAssets(c uipath.Client) ([]uipath.Asset, int, error) {
	var assetList uipath.AssetList

	filters := map[string]string{
		"$top": "1",
	}

	_, err := StoreAsset(c)
	if err != nil {
		return assetList.Value, assetList.Count, err
	}

	aHandler := uipath.AssetHandler{
		Client:   &c,
		FolderId: uint(292388),
	}

	return aHandler.List(filters)
}

func UpdateAsset(c uipath.Client) (uipath.Asset, error) {
	asset, err := StoreAsset(c)
	if err != nil {
		return asset, err
	}

	aHandler := uipath.AssetHandler{
		Client:   &c,
		FolderId: uint(292388),
	}

	updateAsset := uipath.Asset{
		ID:          asset.ID,
		Name:        "NewAssetName",
		ValueType:   "Text",
		StringValue: "Eyyyyyyy",
	}

	return aHandler.Update(updateAsset)
}

func StoreAsset(c uipath.Client) (uipath.Asset, error) {
	aHandler := uipath.AssetHandler{
		Client:   &c,
		FolderId: uint(292388),
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

func GetQueueItemByID(c uipath.Client) (uipath.QueueItem, error) {
	qItem, err := StoreQueueItem(c)
	if err != nil {
		return qItem, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client:   &c,
		FolderId: uint(292388),
	}

	return queueHandler.GetByID(qItem.ID)
}

func ListQueueItems(c uipath.Client) ([]uipath.QueueItem, int, error) {
	var queueItemList uipath.QueueItemList
	filters := map[string]string{
		"$top": "1",
	}

	_, err := StoreQueueItem(c)
	if err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	queueHandler := uipath.QueueItemHandler{
		Client:   &c,
		FolderId: uint(292388),
	}

	return queueHandler.List(filters)
}

func StoreQueueItem(c uipath.Client) (uipath.QueueItem, error) {
	qHandler := uipath.QueueItemHandler{
		Client:   &c,
		FolderId: uint(292388),
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

	return qHandler.Store(qI)
}
