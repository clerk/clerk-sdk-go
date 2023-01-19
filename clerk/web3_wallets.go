package clerk

type Web3Wallet struct {
	ID           string        `json:"id"`
	Object       string        `json:"object"`
	Web3Wallet   string        `json:"web3_wallet"`
	Verification *Verification `json:"verification"`
}
