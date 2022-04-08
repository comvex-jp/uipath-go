package uipath

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RequestError defines how the error looks like from the response
type RequestError struct {
	Message          string `json:"message"`
	ErrorName        string `json:"error"`
	ErrorCode        int    `json:"errorCode"`
	ErrorDescription string `json:"error_description"`
	TraceIDs         string `json:"traceId"`
	ResourceIds      []uint `json:"resourceIds"`
}

func (r *RequestError) Error() string {
	if r.ErrorName != "" {
		return fmt.Sprintf("Request Failed: Error Code(%s) %s", r.ErrorName, r.ErrorDescription)
	}

	return fmt.Sprintf("Request Failed: Error Code(%d) %s", r.ErrorCode, r.Message)
}

//ErrorResponseHandler handles the errors from the uipath response
func ErrorResponseHandler(statusCode int, errResp []byte) error {
	var requestError RequestError

	// Check the response body if it's empty and if it is, we assume that it's an HTTP error.
	if len(errResp) < 1 {
		return errors.New(fmt.Sprintf("HTTP Error %d: %s", statusCode, http.StatusText(statusCode)))
	}

	if err := json.Unmarshal(errResp, &requestError); err != nil {
		return err
	}

	return &requestError
}
