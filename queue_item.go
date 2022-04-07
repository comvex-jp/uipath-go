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
)

// QueueItemHandler struct defines what the asset handler looks like
type QueueItemHandler struct {
	Client *Client
}

// Queue struct defines what the asset model looks like
type QueueItem struct {
	ID                                 uint                   `json:"Id,omitempty"`
	QueueDefinitionID                  uint                   `json:"QueueDefinitionId,omitempty"`
	DeferDate                          string                 `json:"DeferDate"`
	DueDate                            string                 `json:"DueDate"`
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

// QueueItemList defines what the asset model looks like
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
func (q *QueueItemHandler) Store(queueItem QueueItem, folderID uint) (QueueItem, error) {
	var result QueueItem

	queueItemCreateRequest := QueueItemCreateRequest{
		ItemData: queueItem,
	}

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s", q.Client.URLEndpoint, "Queues/UiPathODataSvc.AddQueueItem")

	resp, err := q.Client.SendWithAuthorization("POST", url, queueItemCreateRequest, headers, map[string]string{})
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(resp, &result); err != nil {
		return result, err
	}

	return result, nil
}

// GetByID fetches a queue item by id
func (q *QueueItemHandler) GetByID(ID uint, folderID uint) (QueueItem, error) {
	var queueItem QueueItem

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s(%d)", q.Client.URLEndpoint, "QueueItems", ID)

	resp, err := q.Client.SendWithAuthorization("GET", url, nil, headers, map[string]string{})
	if err != nil {
		return queueItem, err
	}

	if err = json.Unmarshal(resp, &queueItem); err != nil {
		return queueItem, err
	}

	return queueItem, nil
}

// List fetches a list of queue items that can be filtered
func (q *QueueItemHandler) List(filters map[string]string, folderID uint) ([]QueueItem, int, error) {
	var queueItemList QueueItemList

	headers := map[string]string{
		"X-UIPATH-OrganizationUnitId": strconv.Itoa(int(folderID)),
	}

	url := fmt.Sprintf("%s%s", q.Client.URLEndpoint, "QueueItems")

	resp, err := q.Client.SendWithAuthorization("GET", url, nil, headers, filters)
	if err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	if err = json.Unmarshal(resp, &queueItemList); err != nil {
		return queueItemList.Value, queueItemList.Count, err
	}

	return queueItemList.Value, queueItemList.Count, err
}
