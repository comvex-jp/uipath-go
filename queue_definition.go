package uipath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	QueueDefinitionEndpoint = "QueueDefinitions"
)

// QueueDefinitionHandler struct defines what the queuedefinition handler looks like
type QueueDefinitionHandler struct {
	Client            *Client
	QueueDefinitionId uint
	FolderID          uint
}

// QueueDefinition struct defines what the queuedefinition model looks like
type QueueDefinition struct {
	AcceptAutomaticallyRetry           bool      `json:"AcceptAutomaticallyRetry,omitempty"`
	EnforceUniqueReference             bool      `json:"EnforceUniqueReference,omitempty"`
	Encrypted                          bool      `json:"Encrypted,omitempty"`
	IsProcessInCurrentFolder           bool      `json:"IsProcessInCurrentFolder,omitempty"`
	ID                                 uint      `json:"Id"`
	MaxNumberOfRetries                 int32     `json:"MaxNumberOfRetries,omitempty"`
	SlaInMinutes                       int32     `json:"SlaInMinutes,omitempty"`
	RiskSlaInMinutes                   int32     `json:"RiskSlaInMinutes,omitempty"`
	FoldersCount                       int32     `json:"FoldersCount,omitempty"`
	OrganizationUnitId                 int64     `json:"OrganizationUnitId,omitempty"`
	ProcessScheduleId                  int64     `json:"ProcessScheduleId,omitempty"`
	ReleaseId                          int64     `json:"ReleaseId,omitempty"`
	CreationTime                       time.Time `json:"CreationTime,omitempty"`
	Key                                string    `json:"Key"`
	Description                        string    `json:"Description,omitempty"`
	Name                               string    `json:"Name"`
	SpecificDataJsonSchema             string    `json:"SpecificDataJsonSchema,omitempty"`
	OutputDataJsonSchema               string    `json:"OutputDataJsonSchema,omitempty"`
	AnalyticsDataJsonSchema            string    `json:"AnalyticsDataJsonSchema,omitempty"`
	OrganizationUnitFullyQualifiedName string    `json:"OrganizationUnitFullyQualifiedName,omitempty"`
	Tags                               []string  `json:"Tag,omitempty"`
}

// QueueDefinitionList struct defines what response of the list endpoint looks like
type QueueDefinitionList struct {
	Count int               `json:"@odata.count"`
	Value []QueueDefinition `json:"value"`
}

// List fetches a list of queuedefinitions that can be filtered using query parameters
func (q *QueueDefinitionHandler) List(filters map[string]string) ([]QueueDefinition, int, error) {
	var queuedefinitionList QueueDefinitionList

	url := fmt.Sprintf("%s%s", q.Client.BaseURL, QueueDefinitionEndpoint)

	resp, err := q.Client.SendWithAuthorization("GET", url, nil, q.buildHeaders(), filters)
	if err != nil {
		return queuedefinitionList.Value, queuedefinitionList.Count, err
	}

	err = json.Unmarshal(resp, &queuedefinitionList)

	return queuedefinitionList.Value, queuedefinitionList.Count, err
}

func (q *QueueDefinitionHandler) buildHeaders() map[string]string {
	var headers = map[string]string{}

	headers[HeaderOrganizationUnitId] = strconv.Itoa(int(q.FolderID))

	return headers
}
