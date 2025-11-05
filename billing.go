package clerk

// Feature represents a feature associated with a plan.
type Feature struct {
	APIResource

	Object      string `json:"object"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	AvatarURL   string `json:"avatar_url"`
}

// BillingMoney represents money amounts with formatting.
type BillingMoney struct {
	APIResource

	Amount          int64  `json:"amount"`
	AmountFormatted string `json:"amount_formatted"`
	Currency        string `json:"currency"`
	CurrencySymbol  string `json:"currency_symbol"`
}

// BillingProduct represents a product.
type BillingProduct struct {
	APIResource

	Object    string `json:"object"`
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Plans     []Plan `json:"plans"`
}

// Plan represents a billing plan.
type Plan struct {
	APIResource

	Object           string        `json:"object"`
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Fee              *BillingMoney `json:"fee"`
	AnnualMonthlyFee *BillingMoney `json:"annual_monthly_fee"`
	AnnualFee        *BillingMoney `json:"annual_fee"`
	Description      *string       `json:"description"`
	ProductID        string        `json:"product_id"`
	IsDefault        bool          `json:"is_default"`
	IsRecurring      bool          `json:"is_recurring"`
	PubliclyVisible  bool          `json:"publicly_visible"`
	HasBaseFee       bool          `json:"has_base_fee"`
	ForPayerType     string        `json:"for_payer_type"`
	Slug             string        `json:"slug"`
	AvatarURL        *string       `json:"avatar_url"`
	Features         []Feature     `json:"features"`
	FreeTrialEnabled bool          `json:"free_trial_enabled"`
	FreeTrialDays    *int          `json:"free_trial_days"`
}

// PlanList contains a list of plans.
type PlanList struct {
	APIResource

	Data       []Plan `json:"data"`
	TotalCount int64  `json:"total_count"`
}

// BillingPaymentMethod represents a payment method.
type BillingPaymentMethod struct {
	APIResource

	Object                   string  `json:"object"`
	ID                       string  `json:"id"`
	PayerID                  string  `json:"payer_id"`
	PaymentType              string  `json:"payment_type"`
	IsDefault                *bool   `json:"is_default"`
	Gateway                  string  `json:"gateway"`
	GatewayExternalID        string  `json:"gateway_external_id"`
	GatewayExternalAccountID *string `json:"gateway_external_account_id"`
	Last4                    string  `json:"last4"`
	Status                   string  `json:"status"`
	WalletType               *string `json:"wallet_type"`
	CardType                 *string `json:"card_type"`
	ExpiryYear               *int    `json:"expiry_year"`
	ExpiryMonth              *int    `json:"expiry_month"`
	CreatedAt                int64   `json:"created_at"`
	UpdatedAt                int64   `json:"updated_at"`
	IsRemovable              bool    `json:"is_removable"`
}

// BillingSubscriptionItemNextPayment represents next payment info.
type BillingSubscriptionItemNextPayment struct {
	APIResource

	Amount BillingMoney `json:"amount"`
	Date   int64        `json:"date"`
}

// Payer represents a billing payer (user or organization).
type Payer struct {
	APIResource

	Object     string `json:"object"`
	ID         string `json:"id"`
	InstanceID string `json:"instance_id"`

	// User payer only
	UserID    *string `json:"user_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`

	// Org payer only
	OrganizationID   *string `json:"organization_id"`
	OrganizationName *string `json:"organization_name"`

	// Used for both org and user payers
	ImageURL *string `json:"image_url"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// SubscriptionItem represents a billing subscription item.
type SubscriptionItem struct {
	APIResource

	Object          string                              `json:"object"`
	ID              string                              `json:"id"`
	InstanceID      string                              `json:"instance_id"`
	Status          string                              `json:"status"`
	PlanID          *string                             `json:"plan_id"`
	Plan            *Plan                               `json:"plan"`
	PlanPeriod      string                              `json:"plan_period"`
	PaymentMethodID string                              `json:"payment_method_id"`
	PaymentMethod   *BillingPaymentMethod               `json:"payment_method"`
	LifetimePaid    *BillingMoney                       `json:"lifetime_paid"`
	Amount          *BillingMoney                       `json:"amount"`
	NextPayment     *BillingSubscriptionItemNextPayment `json:"next_payment"`
	PayerID         string                              `json:"payer_id"`
	Payer           *Payer                              `json:"payer"`
	IsFreeTrial     bool                                `json:"is_free_trial"`
	PeriodStart     int64                               `json:"period_start"`
	PeriodEnd       *int64                              `json:"period_end"`
	ProrationDate   *string                             `json:"proration_date"`
	CanceledAt      *int64                              `json:"canceled_at"`
	PastDueAt       *int64                              `json:"past_due_at"`
	EndedAt         *int64                              `json:"ended_at"`
	CreatedAt       int64                               `json:"created_at"`
	UpdatedAt       int64                               `json:"updated_at"`
}

// SubscriptionItemList contains a list of subscription items.
type SubscriptionItemList struct {
	APIResource

	Data       []SubscriptionItem `json:"data"`
	TotalCount int64              `json:"total_count"`
}
