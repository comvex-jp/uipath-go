package uipath

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	ValueTypeText       = "Text"
	ValueTypeInteger    = "Integer"
	ValueTypeBool       = "Bool"
	ValueTypeCredential = "Credential"

	ValueScopeGlobal = "Global"

	AssetEndpoint = "Assets"
)

// AssetHandler struct defines what the asset handler looks like
type AssetHandler struct {
	Client   *Client
	FolderId uint
}

// Asset struct defines what the asset model looks like
type Asset struct {
	ID                 uint     `json:"Id"`
	Name               string   `json:"Name"`
	CanBeDeleted       bool     `json:"CanBeDeleted,omitempty"`
	ValueScope         string   `json:"ValueScope,omitempty"`
	ValueType          string   `json:"ValueType"`
	Value              string   `json:"Value,omitempty"`
	StringValue        string   `json:"StringValue,omitempty"`
	BoolValue          bool     `json:"BoolValue,omitempty"`
	IntValue           int      `json:"IntValue,omitempty"`
	CredentialUsername string   `json:"CredentialUsername,omitempty"`
	CredentialPassword string   `json:"CredentialPassword,omitempty"`
	ExternalName       string   `json:"ExternalName,omitempty"`
	CredentialStoreId  int      `json:"CredentialStoreId,omitempty"`
	HasDefaultValue    bool     `json:"HasDefault,omitempty"`
	Description        string   `json:"Description,omitempty"`
	FolderCount        int      `json:"FolderCount,omitempty"`
	Tags               []string `json:"Tag,omitempty"`
	KeyValueList       []string `json:"KeyValueList,omitempty"`
}

// AssetList struct defines what the asset list looks like
type AssetList struct {
	Count int     `json:"@odata.count"`
	Value []Asset `json:"value"`
}

// Tag struct defines what the tag item looks like
type Tag struct {
	Name         string `json:"Name"`
	DisplayName  string `json:"DisplayName"`
	Value        string `json:"Value"`
	DisplayValue string `json:"DisplayValue,omitempty"`
}

// GetByID fetches the asset by id
func (a *AssetHandler) GetByID(ID uint) (Asset, error) {
	var asset Asset

	url := fmt.Sprintf("%s%s(%d)", a.Client.BaseURL, AssetEndpoint, ID)

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, a.buildHeaders(), map[string]string{})
	if err != nil {
		return asset, err
	}

	err = json.Unmarshal(resp, &asset)

	return asset, err
}

// GetByName fetches the asset by name
func (a *AssetHandler) GetByName(name string) (Asset, error) {
	var assetList AssetList
	var asset Asset

	url := fmt.Sprintf("%s%s?$filter=Name eq '%s'", a.Client.BaseURL, AssetEndpoint, name)

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, a.buildHeaders(), map[string]string{})
	if err != nil {
		return asset, err
	}

	if err = json.Unmarshal(resp, &assetList); err != nil {
		return asset, err
	}

	if assetList.Count < 1 {
		return asset, nil
	}

	return assetList.Value[0], nil
}

// List fetches a list of assets that can be filtered using query parameters
func (a *AssetHandler) List(filters map[string]string) ([]Asset, int, error) {
	var assetList AssetList

	url := fmt.Sprintf("%s%s", a.Client.BaseURL, AssetEndpoint)

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, a.buildHeaders(), filters)
	if err != nil {
		return assetList.Value, assetList.Count, err
	}

	err = json.Unmarshal(resp, &assetList)

	return assetList.Value, assetList.Count, err
}

// Store creates and saves an asset on the orchestrator
func (a *AssetHandler) Store(asset Asset) (Asset, error) {
	var result Asset

	url := fmt.Sprintf("%s%s", a.Client.BaseURL, AssetEndpoint)

	resp, err := a.Client.SendWithAuthorization("POST", url, asset, a.buildHeaders(), map[string]string{})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(resp, &result)

	return result, err
}

// Update updates an asset
func (a *AssetHandler) Update(asset Asset) (Asset, error) {
	url := fmt.Sprintf("%s%s(%d)", a.Client.BaseURL, AssetEndpoint, asset.ID)

	_, err := a.Client.SendWithAuthorization("PUT", url, asset, a.buildHeaders(), map[string]string{})
	if err != nil {
		return asset, err
	}

	return a.GetByID(asset.ID)
}

func (a *AssetHandler) buildHeaders() map[string]string {
	var headers = map[string]string{}

	headers[HeaderOrganizationUnitId] = strconv.Itoa(int(a.FolderId))

	return headers
}
