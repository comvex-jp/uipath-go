package uipath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	PriorityLow    = "Low"
	PriorityNormal = "Normal"
	PriorityHigh   = "High"

	QueueItemEndpoint    = "QueueItem"
	QueueAddItemEndpoint = "Queues/UiPathODataSvc.AddQueueItem"
)

// QueueItemHandler struct defines what the queue item handler looks like
type QueueItemHandler struct {
	Client   *Client
	FolderId uint
}

// QueueItem struct defines what the queue item model looks like
type QueueItem struct {
	ID                                 uint                   `json:"Id,omitempty"`
	QueueDefinitionID                  uint                   `json:"QueueDefinitionId,omitempty"`
	DeferDate                          string                 `json:"DeferDate,omitempty"`
	DueDate                            string                 `json:"DueDate,omitempty"`
	RiskSlaDate                        string                 `json:"RiskSlaDate,omitempty"`
	Priority                           string                 `json:"Priority"`
	Status                             string                 `json:"Status,omitempty"`
	ReviewStatus                       string                 `json:"ReviewStatus,omitempty"`
	ReviewerUserID                     uint                   `json:"ReviewerUserId,omitempty"`
	Key                                string                 `json:"Key,omitempty"`
	Reference                          string                 `json:"Reference,omitempty"`
	ProcessingExceptionType            string                 `json:"ProcessingExceptionType,omitempty"`
	StartProcessing                    string                 `json:"StartProcessing,omitempty"`
	EndProcessing                      string                 `json:"EndProcessing,omitempty"`
	SecondsInPreviousAttempts          uint                   `json:"SecondsInPreviousAttempts,omitempty"`
	AncestorID                         *uint                  `json:"AncestorId,omitempty"`
	RetryNumber                        uint                   `json:"RetryNumber,omitempty"`
	SpecificData                       string                 `json:"SpecificData,omitempty"`
	CreationTime                       string                 `json:"CreationTime,omitempty"`
	Progress                           string                 `json:"Progress,omitempty"`
	RowVersion                         string                 `json:"RowVersion,omitempty"`
	OrganizationUnitID                 uint                   `json:"OrganizationUnitId,omitempty"`
	OrganizationUnitFullyQualifiedName string                 `json:"OrganizationUnitFullyQualifiedName,omitempty"`
	ProcessingException                *ProcessingException   `json:"ProcessingException,omitempty"`
	SpecificContent                    map[string]interface{} `json:"SpecificContent,omitempty"`
	Output                             map[string]interface{} `json:"Output,omitempty"`
	Name                               string                 `json:"Name"`
}

// QueueItemList defines what the queue item list model looks like
type QueueItemList struct {
	Count int         `json:"@odata.count"`
	Value []QueueItem `json:"value"`
}

// QueueItemCreateRequest defines how the request looks like when creating a queue item
type QueueItemCreateRequest struct {
	ItemData QueueItem `json:"itemData"`
}

// ProcessingException defines the structure of the queue item exception
type ProcessingException struct {
	Reason                  string
	Details                 string
	Type                    string
	AssociatedImageFilePath string
	CreationTime            time.Time
}

// Store creates and stores a queue item in the uipath orchestrator
func (q *QueueItemHandler) Store(queueItem QueueItem) (QueueItem, error) {
	var result QueueItem

	queueItemCreateRequest := QueueItemCreateRequest{
		ItemData: queueItem,
	}

	url := fmt.Sprintf("%s%s", q.Client.BaseURL, QueueAddItemEndpoint)

	resp, err := q.Client.SendWithAuthorization("POST", url, queueItemCreateRequest, q.buildHeaders(), map[string]string{})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(resp, &result)

	return result, err
}

// GetByID fetches a queue item by id
func (q *QueueItemHandler) GetByID(ID uint) (QueueItem, error) {
	var queueItem QueueItem

	url := fmt.Sprintf("%s%s(%d)", q.Client.BaseURL, QueueItemEndpoint, ID)

	resp, err := q.Client.SendWithAuthorization("GET", url, nil, q.buildHeaders(), map[string]string{})
	if err != nil {
		return queueItem, err
	}

	err = json.Unmarshal(resp, &queueItem)

	return queueItem, err
}

// List fetches a list of queue items that can be filtered
func (q *QueueItemHandler) List(filters map[string]string) ([]QueueItem, int, error) {
	var queueItemList QueueItemList

	url := fmt.Sprintf("%s%s", q.Client.BaseURL, QueueItemEndpoint)

	resp, err := q.Client.SendWithAuthorization("GET", url, nil, q.buildHeaders(), filters)
	if err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	err = json.Unmarshal(resp, &queueItemList)

	return queueItemList.Value, queueItemList.Count, err
}

func (q *QueueItemHandler) buildHeaders() map[string]string {
	var headers = map[string]string{}

	headers[HeaderOrganizationUnitId] = strconv.Itoa(int(q.FolderId))

	return headers
}
