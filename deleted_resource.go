package clerk

// DeletedResource describes an API resource that is no longer
// available.
// It's usually encountered as a result of delete API operations.
type DeletedResource struct {
	APIResource
	ID      string `json:"id,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
