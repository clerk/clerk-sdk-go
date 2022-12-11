package clerk

import "fmt"

type SignUpService service

type SignUp struct {
	Object           string              `json:"object"`
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	RequiredFields   []string            `json:"required_fields"`
	OptionalFields   []string            `json:"optional_fields"`
	MissingFields    []string            `json:"missing_fields"`
	UnverifiedFields []string            `json:"unverified_fields"`
	Verifications    SignUpVerifications `json:"verifications"`

	Username     *string `json:"username"`
	EmailAddress *string `json:"email_address"`
	PhoneNumber  *string `json:"phone_number"`
	Web3Wallet   *string `json:"web3_wallet"`

	PasswordEnabled  bool        `json:"password_enabled"`
	FirstName        *string     `json:"first_name"`
	LastName         *string     `json:"last_name"`
	UnsafeMetadata   interface{} `json:"unsafe_metadata,omitempty"`
	PublicMetadata   interface{} `json:"public_metadata,omitempty"`
	CustomAction     bool        `json:"custom_action"`
	ExternalID       *string     `json:"external_id"`
	CreatedSessionID *string     `json:"created_session_id"`
	CreatedUserID    *string     `json:"created_user_id"`
	AbandonAt        int64       `json:"abandon_at"`

	IdentificationRequirements  [][]string    `json:"identification_requirements"`   // DX: Deprecated
	MissingRequirements         []string      `json:"missing_requirements"`          // DX: Deprecated
	EmailAddressVerification    *Verification `json:"email_address_verification"`    // DX: Deprecated
	PhoneNumberVerification     *Verification `json:"phone_number_verification"`     // DX: Deprecated
	ExternalAccountStrategy     *string       `json:"external_account_strategy"`     // DX: Deprecated
	ExternalAccountVerification *Verification `json:"external_account_verification"` // DX: Deprecated
	ExternalAccount             interface{}   `json:"external_account,omitempty"`    // DX: Deprecated
}

type SignUpVerifications struct {
	EmailAddress    *signUpVerification `json:"email_address"`
	PhoneNumber     *signUpVerification `json:"phone_number"`
	Web3Wallet      *signUpVerification `json:"web3_wallet"`
	ExternalAccount *Verification       `json:"external_account"`
}

type signUpVerification struct {
	*Verification
	NextAction          string   `json:"next_action"`
	SupportedStrategies []string `json:"supported_strategies"`
}

func (s *SignUpService) Read(signUpID string) (*SignUp, error) {
	signUpURL := fmt.Sprintf("%s/%s", SignUpsURL, signUpID)
	req, _ := s.client.NewRequest("GET", signUpURL)

	var signUp SignUp

	_, err := s.client.Do(req, &signUp)
	if err != nil {
		return nil, err
	}

	return &signUp, nil
}

// FIXME how to define "not set" semantics for external_id?
type UpdateSignUp struct {
	CustomAction *bool   `json:"custom_action" form:"custom_action"`
	ExternalID   *string `json:"external_id" form:"external_id"`
}

func (s *SignUpService) Update(signUpID string, updateRequest *UpdateSignUp) (*SignUp, error) {
	signUpURL := fmt.Sprintf("%s/%s", SignUpsURL, signUpID)
	req, _ := s.client.NewRequest("PATCH", signUpURL, updateRequest)

	var signUp SignUp

	_, err := s.client.Do(req, &signUp)
	if err != nil {
		return nil, err
	}

	return &signUp, nil
}
