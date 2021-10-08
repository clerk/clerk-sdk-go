package clerk

type DeleteResponse struct {
	ID      string `json:"id,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
