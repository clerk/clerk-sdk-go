package clerk

import (
	"net/http"
)

type InstanceService service

type UpdateInstanceParams struct {
	// TestMode can be used to toggle test mode for this instance.
	// Defaults to true for development instances.
	TestMode *bool `json:"test_mode,omitempty"`

	// HIBP is used to configure whether Clerk should use the
	// "Have I Been Pawned" service to check passwords against
	// known security breaches.
	// By default, this is enabled in all instances.
	HIBP *bool `json:"hibp,omitempty"`
}

func (s *InstanceService) Update(params UpdateInstanceParams) error {
	req, _ := s.client.NewRequest(http.MethodPatch, "instance", &params)

	_, err := s.client.Do(req, nil)
	return err
}
