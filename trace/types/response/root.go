package response

import "net/http"

type Response struct {
	Status int64 `json:"status"`
}

func (r Response) IsCallFailed() bool {
	statusCode := r.Status
	return statusCode == http.StatusBadRequest || (statusCode >= http.StatusInternalServerError && statusCode < 600)
}

func (r Response) IsCallSuccessful() bool {
	statusCode := r.Status
	return statusCode <= http.StatusMultipleChoices || statusCode == http.StatusNotModified || statusCode == http.StatusTemporaryRedirect
}

func (r Response) IsCallBlocked() bool {
	statusCode := r.Status
	return statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden || statusCode == http.StatusTooManyRequests
}
