package clerk

type SessionToken struct {
	APIResource
	Object string `json:"object"`
	JWT    string `json:"jwt"`
}
