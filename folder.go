package uipath

import (
	"encoding/json"
	"fmt"
)

const (
	FolderEndpoint = "Folders"
)

// FolderHandler struct defines what the folder handler looks like
type FolderHandler struct {
	Client *Client
}

// Folder struct defines what the folder model looks like
type Folder struct {
	IsPersonal         bool   `json:"IsPersonal,omitempty"`
	ID                 uint   `json:"Id"`
	ParentId           int64  `json:"ParentId,omitempty"`
	Key                string `json:"Key"`
	DisplayName        string `json:"DisplayName,omitempty"`
	FullyQualifiedName string `json:"FullyQualifiedName"`
	Description        string `json:"Description,omitempty"`
	ProvisionType      string `json:"ProvisionType,omitempty"`
	PermissionModel    string `json:"PermissionModel,omitempty"`
	ParentKey          string `json:"ParentKey,omitempty"`
	FeedType           string `json:"FeedType,omitempty"`
}

// FolderList struct defines what response of the list endpoint looks like
type FolderList struct {
	Count int      `json:"@odata.count"`
	Value []Folder `json:"value"`
}

// List fetches a list of folders that can be filtered using query parameters
func (f *FolderHandler) List(filters map[string]string) ([]Folder, int, error) {
	var folderList FolderList

	url := fmt.Sprintf("%s%s", f.Client.BaseURL, FolderEndpoint)

	resp, err := f.Client.SendWithAuthorization("GET", url, nil, nil, filters)
	if err != nil {
		return folderList.Value, folderList.Count, err
	}

	err = json.Unmarshal(resp, &folderList)

	return folderList.Value, folderList.Count, err
}
