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
	ValueTypeCredential = "Credentials"

	ValueScopeGlobal = "Global"
)

// AssetHandler struct defines what the asset handler looks like
type AssetHandler struct {
	Client *Client
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
	Value []Asset `json:"Value"`
}

// Tag struct defines what the tag item looks like
type Tag struct {
	Name         string `json:"Name"`
	DisplayName  string `json:"DisplayName"`
	Value        string `json:"Value"`
	DisplayValue string `json:"DisplayValue,omitempty"`
}

// GetByID fetches the asset by id
func (a *AssetHandler) GetByID(ID uint, folderID uint) (Asset, error) {
	var asset Asset

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s(%d)", a.Client.URLEndpoint, "Assets", ID)

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, headers, map[string]string{})
	if err != nil {
		return asset, err
	}

	if err = json.Unmarshal(resp, &asset); err != nil {
		return asset, err
	}

	return asset, nil
}

// GetByID fetches the asset by name
func (a *AssetHandler) GetByName(name string, folderID uint) (Asset, error) {
	var asset Asset

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s?$filter=Name eq '%s'", a.Client.URLEndpoint, "Assets", name)

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, headers, map[string]string{})
	if err != nil {
		return asset, err
	}

	if err = json.Unmarshal(resp, &asset); err != nil {
		return asset, err
	}

	return asset, nil
}

// List fetches a list of assets that can be filtered using query parameters
func (a *AssetHandler) List(filters map[string]string, folderID uint) ([]Asset, int, error) {
	var assetList AssetList

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s", a.Client.URLEndpoint, "Assets")

	resp, err := a.Client.SendWithAuthorization("GET", url, nil, headers, filters)
	if err != nil {
		return assetList.Value, assetList.Count, err
	}

	if err = json.Unmarshal(resp, &assetList); err != nil {
		return assetList.Value, assetList.Count, err
	}

	return assetList.Value, assetList.Count, nil
}

// Store creates and saves an asset on the orchestrator
func (a *AssetHandler) Store(asset Asset, folderID uint) (Asset, error) {
	var result Asset

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s", a.Client.URLEndpoint, "Assets")

	resp, err := a.Client.SendWithAuthorization("POST", url, asset, headers, map[string]string{})
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(resp, &result); err != nil {
		return result, err
	}

	return result, nil
}

// Update updates an asset
func (a *AssetHandler) Update(asset Asset, folderID uint) (Asset, error) {
	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s(%d)", a.Client.URLEndpoint, "Assets", asset.ID)

	_, err := a.Client.SendWithAuthorization("PUT", url, asset, headers, map[string]string{})
	if err != nil {
		return asset, err
	}

	updatedAsset, err := a.GetByID(asset.ID, folderID)
	if err != nil {
		return asset, err
	}

	return updatedAsset, nil
}
