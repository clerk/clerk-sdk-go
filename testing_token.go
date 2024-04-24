package clerk

type TestingToken struct {
	APIResource
	Object    string `json:"object"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
