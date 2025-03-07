package http

import (
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

type ErrorPayload struct {
	Type     string         `json:"type"`
	Reason   string         `json:"reason"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type ClerkErrorResponse struct {
	ClerkError ErrorPayload `json:"clerk_error"`
}

func NewSessionReverificationErrorPayload(missingConfig clerk.SessionReverificationPolicy) ClerkErrorResponse {
	return ClerkErrorResponse{
		ClerkError: ErrorPayload{
			Type:   "forbidden",
			Reason: "reverification-error",
			Metadata: map[string]any{
				"reverification": missingConfig,
			},
		},
	}
}

func WriteNeedsReverificationResponse(w http.ResponseWriter, missingConfig clerk.SessionReverificationPolicy) {
	payload := NewSessionReverificationErrorPayload(missingConfig)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(payload)
}
