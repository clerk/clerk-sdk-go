package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/stretchr/testify/require"
)

func TestBillingClientListPlans(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"plan","id":"plan_123","name":"Basic Plan","for_payer_type":"user","features":[{"object":"feature","id":"feature_456","name":"Feature 1","key":"feature_1"}]}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/plans",
			Query: &url.Values{
				"limit":      []string{"10"},
				"offset":     []string{"0"},
				"payer_type": []string{"user"},
			},
		},
	}
	client := NewClient(config)
	params := &ListPlansParams{
		PayerType: clerk.String("user"),
	}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)
	planList, err := client.ListPlans(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), planList.TotalCount)
	require.Equal(t, 1, len(planList.Data))
	require.Equal(t, "plan_123", planList.Data[0].ID)
	require.Equal(t, "Basic Plan", planList.Data[0].Name)
	require.Equal(t, "user", planList.Data[0].ForPayerType)
	require.Equal(t, 1, len(planList.Data[0].Features))
	require.Equal(t, "feature_456", planList.Data[0].Features[0].ID)
	require.Equal(t, "Feature 1", planList.Data[0].Features[0].Name)
}

func TestBillingClientListPlans_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"list-plans-error-code"
	}],
	"clerk_trace_id":"list-plans-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.ListPlans(context.Background(), &ListPlansParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "list-plans-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "list-plans-error-code", apiErr.Errors[0].Code)
}

func TestBillingClientListSubscriptionItems(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","payer_id":"payer_456","plan_id":"plan_789","status":"active","period_start":1640995200,"period_end":1643673600,"payer":{"object":"payer","id":"payer_456","user_id":"user_456","first_name":"John","last_name":"Doe","email":"john@example.com","created_at":1640995200,"updated_at":1640995200},"plan":{"object":"plan","id":"plan_789","name":"Pro Plan","payer_type":["user"],"features":[]},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
			Query: &url.Values{
				"limit":   []string{"10"},
				"offset":  []string{"0"},
				"user_id": []string{"user_456"},
				"status":  []string{"active"},
			},
		},
	}
	client := NewClient(config)
	params := &ListSubscriptionItemsParams{
		UserID: clerk.String("user_456"),
		Status: clerk.String("active"),
	}
	params.Limit = clerk.Int64(10)
	params.Offset = clerk.Int64(0)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), subscriptionItemList.TotalCount)
	require.Equal(t, 1, len(subscriptionItemList.Data))
	require.Equal(t, "sub_item_123", subscriptionItemList.Data[0].ID)
	require.Equal(t, "payer_456", subscriptionItemList.Data[0].PayerID)
	require.Equal(t, "plan_789", *subscriptionItemList.Data[0].PlanID)
	require.Equal(t, "active", subscriptionItemList.Data[0].Status)
	require.Equal(t, int64(1640995200), subscriptionItemList.Data[0].PeriodStart)
	require.Equal(t, int64(1643673600), *subscriptionItemList.Data[0].PeriodEnd)
	require.NotNil(t, subscriptionItemList.Data[0].Payer)
	require.Equal(t, "payer_456", subscriptionItemList.Data[0].Payer.ID)
	require.Equal(t, "John", *subscriptionItemList.Data[0].Payer.FirstName)
	require.NotNil(t, subscriptionItemList.Data[0].Plan)
	require.Equal(t, "plan_789", subscriptionItemList.Data[0].Plan.ID)
	require.Equal(t, "Pro Plan", subscriptionItemList.Data[0].Plan.Name)
}

func TestBillingClientListSubscriptionItems_Error(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"list-subscription-items-error-code"
	}],
	"clerk_trace_id":"list-subscription-items-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "list-subscription-items-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "list-subscription-items-error-code", apiErr.Errors[0].Code)
}

func TestBillingClientCancelSubscriptionItem(t *testing.T) {
	t.Parallel()
	subscriptionItemID := "sub_item_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"subscription_item","id":"%s","payer_id":"payer_456","plan_id":"plan_789","status":"canceled","period_start":1640995200,"period_end":1643673600,"created_at":1640995200,"updated_at":1640995260}`, subscriptionItemID)),
			Method: http.MethodDelete,
			Path:   "/v1/billing/subscription_items/" + subscriptionItemID,
			Query: &url.Values{
				"end_now": []string{"true"},
			},
		},
	}
	client := NewClient(config)
	subscriptionItem, err := client.CancelSubscriptionItem(context.Background(), subscriptionItemID, &CancelSubscriptionItemParams{
		EndNow: clerk.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, subscriptionItemID, subscriptionItem.ID)
	require.Equal(t, "canceled", subscriptionItem.Status)
	require.Equal(t, int64(1640995260), subscriptionItem.UpdatedAt)
}

func TestBillingClientCancelSubscriptionItem_Error(t *testing.T) {
	t.Parallel()
	subscriptionItemID := "sub_item_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusNotFound,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"subscription-item-not-found"
	}],
	"clerk_trace_id":"cancel-subscription-item-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.CancelSubscriptionItem(context.Background(), subscriptionItemID, &CancelSubscriptionItemParams{})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "cancel-subscription-item-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "subscription-item-not-found", apiErr.Errors[0].Code)
}

func TestBillingClientExtendFreeTrial(t *testing.T) {
	t.Parallel()
	subscriptionItemID := "sub_item_123"
	extendTo := "2024-12-31T23:59:59Z"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			In:     json.RawMessage(fmt.Sprintf(`{"extend_to":"%s"}`, extendTo)),
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"subscription_item","id":"%s","payer_id":"payer_456","plan_id":"plan_789","status":"trialing","period_start":1640995200,"period_end":1735689599,"created_at":1640995200,"updated_at":1640995260}`, subscriptionItemID)),
			Method: http.MethodPost,
			Path:   "/v1/billing/subscription_items/" + subscriptionItemID + "/extend_free_trial",
		},
	}
	client := NewClient(config)
	subscriptionItem, err := client.ExtendFreeTrial(context.Background(), subscriptionItemID, &ExtendFreeTrialParams{
		ExtendTo: extendTo,
	})
	require.NoError(t, err)
	require.Equal(t, subscriptionItemID, subscriptionItem.ID)
	require.Equal(t, "trialing", subscriptionItem.Status)
	require.Equal(t, int64(1735689599), *subscriptionItem.PeriodEnd)
	require.Equal(t, int64(1640995260), subscriptionItem.UpdatedAt)
}

func TestBillingClientExtendFreeTrial_Error(t *testing.T) {
	t.Parallel()
	subscriptionItemID := "sub_item_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Status: http.StatusBadRequest,
			Out: json.RawMessage(`{
  "errors":[{
		"code":"subscription-item-not-in-free-trial"
	}],
	"clerk_trace_id":"extend-free-trial-trace-id"
}`),
		},
	}
	client := NewClient(config)
	_, err := client.ExtendFreeTrial(context.Background(), subscriptionItemID, &ExtendFreeTrialParams{
		ExtendTo: "2024-12-31T23:59:59Z",
	})
	require.Error(t, err)
	apiErr, ok := err.(*clerk.APIErrorResponse)
	require.True(t, ok)
	require.Equal(t, "extend-free-trial-trace-id", apiErr.TraceID)
	require.Equal(t, 1, len(apiErr.Errors))
	require.Equal(t, "subscription-item-not-in-free-trial", apiErr.Errors[0].Code)
}

func TestBillingClientCancelSubscriptionItemWithoutEndNow(t *testing.T) {
	t.Parallel()
	subscriptionItemID := "sub_item_123"
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(fmt.Sprintf(`{"object":"subscription_item","id":"%s","payer_id":"payer_456","plan_id":"plan_789","status":"canceled","period_start":1640995200,"period_end":1643673600,"created_at":1640995200,"updated_at":1640995260}`, subscriptionItemID)),
			Method: http.MethodDelete,
			Path:   "/v1/billing/subscription_items/" + subscriptionItemID,
			Query: &url.Values{
				"end_now": []string{"false"},
			},
		},
	}
	client := NewClient(config)
	subscriptionItem, err := client.CancelSubscriptionItem(context.Background(), subscriptionItemID, &CancelSubscriptionItemParams{
		EndNow: clerk.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, subscriptionItemID, subscriptionItem.ID)
	require.Equal(t, "canceled", subscriptionItem.Status)
}

func TestBillingClientListPlansWithoutFilters(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"plan","id":"plan_123","name":"Basic Plan","for_payer_type":"user","features":[]},{"object":"plan","id":"plan_456","name":"Pro Plan","for_payer_type":"organization","features":[]}],"total_count": 2}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/plans",
		},
	}
	client := NewClient(config)
	planList, err := client.ListPlans(context.Background(), &ListPlansParams{})
	require.NoError(t, err)
	require.Equal(t, int64(2), planList.TotalCount)
	require.Equal(t, 2, len(planList.Data))
	require.Equal(t, "plan_123", planList.Data[0].ID)
	require.Equal(t, "Basic Plan", planList.Data[0].Name)
	require.Equal(t, "user", planList.Data[0].ForPayerType)
	require.Equal(t, "plan_456", planList.Data[1].ID)
	require.Equal(t, "Pro Plan", planList.Data[1].Name)
	require.Equal(t, "organization", planList.Data[1].ForPayerType)
}

func TestBillingClientListSubscriptionItemsWithMultipleFilters(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","payer_id":"payer_456","plan_id":"plan_789","status":"active","period_start":1640995200,"period_end":1643673600,"payer":{"object":"payer","id":"payer_456","user_id":"user_456","first_name":"John","last_name":"Doe","email":"john@example.com","created_at":1640995200,"updated_at":1640995200},"plan":{"object":"plan","id":"plan_789","name":"Pro Plan","payer_type":["user"],"features":[]},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
			Query: &url.Values{
				"user_id":      []string{"user_456"},
				"plan_id":      []string{"plan_789"},
				"status":       []string{"active"},
				"payer_type":   []string{"user"},
				"include_free": []string{"false"},
			},
		},
	}
	client := NewClient(config)
	params := &ListSubscriptionItemsParams{
		UserID:      clerk.String("user_456"),
		PlanID:      clerk.String("plan_789"),
		Status:      clerk.String("active"),
		PayerType:   clerk.String("user"),
		IncludeFree: clerk.Bool(false),
	}
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), subscriptionItemList.TotalCount)
	require.Equal(t, 1, len(subscriptionItemList.Data))
	require.Equal(t, "sub_item_123", subscriptionItemList.Data[0].ID)
	require.Equal(t, "payer_456", subscriptionItemList.Data[0].PayerID)
	require.Equal(t, "plan_789", *subscriptionItemList.Data[0].PlanID)
	require.Equal(t, "active", subscriptionItemList.Data[0].Status)
}

func TestBillingClientListPlans_WithoutOptionalFields(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"plan","id":"plan_123","name":"Basic Plan","product_id":"prod_123","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":false,"slug":"basic","free_trial_enabled":false}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/plans",
		},
	}
	client := NewClient(config)
	planList, err := client.ListPlans(context.Background(), &ListPlansParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), planList.TotalCount)
	require.Equal(t, 1, len(planList.Data))

	plan := planList.Data[0]
	require.Equal(t, "plan_123", plan.ID)
	require.Equal(t, "Basic Plan", plan.Name)
	require.Equal(t, "user", plan.ForPayerType)

	// These fields have omitempty and are absent from the JSON response
	require.Empty(t, plan.Features, "Features should be empty when omitted from response")
	require.Nil(t, plan.Description, "Description should be nil when omitted from response")
	require.Nil(t, plan.AvatarURL, "AvatarURL should be nil when omitted from response")
}

func TestBillingProduct_WithoutOptionalPlans(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"product","id":"prod_123","slug":"basic-product","currency":"USD","name":"Basic Product","is_default":false}`
	var product clerk.BillingProduct
	err := json.Unmarshal([]byte(jsonData), &product)
	require.NoError(t, err)
	require.Equal(t, "prod_123", product.ID)
	require.Equal(t, "Basic Product", product.Name)
	require.Empty(t, product.Plans, "Plans should be empty when omitted from response")
}

func TestSubscriptionItem_WithoutOptionalNestedObjects(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","instance_id":"ins_123","status":"active","payer_id":"payer_456","plan_id":"plan_789","plan_period":"monthly","payment_method_id":"pm_123","is_free_trial":false,"period_start":1640995200,"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
		},
	}
	client := NewClient(config)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.NoError(t, err)
	require.Equal(t, int64(1), subscriptionItemList.TotalCount)

	subItem := subscriptionItemList.Data[0]
	require.Equal(t, "sub_item_123", subItem.ID)
	require.Equal(t, "active", subItem.Status)

	// These fields have omitempty and are absent from the JSON response
	require.Nil(t, subItem.Plan, "Plan should be nil when omitted from response")
	require.Nil(t, subItem.PaymentMethod, "PaymentMethod should be nil when omitted from response")
	// Payer is non-nullable now, so it should have zero values but not nil
	require.Equal(t, "", subItem.Payer.ID, "Payer should have zero values when omitted")
}

func TestPayer_WithoutOptionalFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payer","id":"payer_123","instance_id":"ins_123","user_id":"user_456"}`
	var payer clerk.Payer
	err := json.Unmarshal([]byte(jsonData), &payer)
	require.NoError(t, err)
	require.Equal(t, "payer_123", payer.ID)
	require.Equal(t, "user_456", *payer.UserID)

	// These fields have omitempty and are absent from the JSON response
	require.Nil(t, payer.FirstName, "FirstName should be nil when omitted")
	require.Nil(t, payer.LastName, "LastName should be nil when omitted")
	require.Nil(t, payer.Email, "Email should be nil when omitted")
	require.Equal(t, "", payer.ImageURL, "ImageURL should be empty string when omitted")
	require.Equal(t, int64(0), payer.CreatedAt, "CreatedAt should be zero when omitted")
	require.Equal(t, int64(0), payer.UpdatedAt, "UpdatedAt should be zero when omitted")
}

func TestBillingPaymentMethod_WithoutOptionalCardFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"bank_account","gateway":"stripe","gateway_external_id":"ext_123","last4":"1234","status":"active","is_removable":true}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonData), &paymentMethod)
	require.NoError(t, err)
	require.Equal(t, "pm_123", paymentMethod.ID)
	require.Equal(t, "bank_account", paymentMethod.PaymentType)

	// These fields have omitempty and are absent from the JSON response
	require.Nil(t, paymentMethod.WalletType, "WalletType should be nil when omitted")
	require.Nil(t, paymentMethod.CardType, "CardType should be nil when omitted")
	require.Nil(t, paymentMethod.ExpiryYear, "ExpiryYear should be nil when omitted")
	require.Nil(t, paymentMethod.ExpiryMonth, "ExpiryMonth should be nil when omitted")
	require.Equal(t, int64(0), paymentMethod.CreatedAt, "CreatedAt should be zero when omitted")
	require.Equal(t, int64(0), paymentMethod.UpdatedAt, "UpdatedAt should be zero when omitted")

	// IsRemovable is now non-nullable, so it should have a value
	require.True(t, paymentMethod.IsRemovable)
}

func TestBillingPaymentMethod_NonNullableIsRemovable(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_123","last4":"4242","status":"active","is_removable":false}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonData), &paymentMethod)
	require.NoError(t, err)
	require.Equal(t, "pm_123", paymentMethod.ID)
	require.False(t, paymentMethod.IsRemovable, "IsRemovable should be false as a non-nullable bool")

	// Test with true value
	jsonDataTrue := `{"object":"payment_method","id":"pm_456","payer_id":"payer_789","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_456","last4":"5555","status":"active","is_removable":true}`
	err = json.Unmarshal([]byte(jsonDataTrue), &paymentMethod)
	require.NoError(t, err)
	require.True(t, paymentMethod.IsRemovable, "IsRemovable should be true as a non-nullable bool")
}

func TestBillingSubscriptionItemNextPayment_NonNullableFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"amount":{"amount":2999,"amount_formatted":"$29.99","currency":"USD","currency_symbol":"$"},"date":1672531200}`
	var nextPayment clerk.BillingSubscriptionItemNextPayment
	err := json.Unmarshal([]byte(jsonData), &nextPayment)
	require.NoError(t, err)

	require.Equal(t, int64(2999), nextPayment.Amount.Amount)
	require.Equal(t, "$29.99", nextPayment.Amount.AmountFormatted)
	require.Equal(t, "USD", nextPayment.Amount.Currency)
	require.Equal(t, "$", nextPayment.Amount.CurrencySymbol)

	require.Equal(t, int64(1672531200), nextPayment.Date)
}

func TestSubscriptionItem_NonNullablePayer(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","instance_id":"ins_123","status":"active","payer_id":"payer_456","plan_id":"plan_789","plan_period":"monthly","payment_method_id":"pm_123","is_free_trial":false,"period_start":1640995200,"payer":{"object":"payer","id":"payer_456","instance_id":"ins_123","user_id":"user_456","first_name":"Jane","last_name":"Smith","email":"jane@example.com"},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
		},
	}
	client := NewClient(config)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.NoError(t, err)

	subItem := subscriptionItemList.Data[0]

	require.Equal(t, "payer_456", subItem.Payer.ID)
	require.Equal(t, "Jane", *subItem.Payer.FirstName)
	require.Equal(t, "Smith", *subItem.Payer.LastName)
	require.Equal(t, "jane@example.com", *subItem.Payer.Email)
}

func TestSubscriptionItem_NonNullableAmountAndLifetimePaid(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","instance_id":"ins_123","status":"active","payer_id":"payer_456","plan_id":"plan_789","plan_period":"monthly","payment_method_id":"pm_123","is_free_trial":false,"period_start":1640995200,"amount":{"amount":4999,"amount_formatted":"$49.99","currency":"USD","currency_symbol":"$"},"lifetime_paid":{"amount":14997,"amount_formatted":"$149.97","currency":"USD","currency_symbol":"$"},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
		},
	}
	client := NewClient(config)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.NoError(t, err)

	subItem := subscriptionItemList.Data[0]

	require.Equal(t, int64(4999), subItem.Amount.Amount)
	require.Equal(t, "$49.99", subItem.Amount.AmountFormatted)
	require.Equal(t, "USD", subItem.Amount.Currency)
	require.Equal(t, int64(14997), subItem.LifetimePaid.Amount)
	require.Equal(t, "$149.97", subItem.LifetimePaid.AmountFormatted)
	require.Equal(t, "USD", subItem.LifetimePaid.Currency)
}

func TestSubscriptionItem_NonNullableNextPayment(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","instance_id":"ins_123","status":"active","payer_id":"payer_456","plan_id":"plan_789","plan_period":"monthly","payment_method_id":"pm_123","is_free_trial":false,"period_start":1640995200,"next_payment":{"amount":{"amount":1999,"amount_formatted":"$19.99","currency":"USD","currency_symbol":"$"},"date":1672531200},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
		},
	}
	client := NewClient(config)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.NoError(t, err)

	subItem := subscriptionItemList.Data[0]

	require.Equal(t, int64(1999), subItem.NextPayment.Amount.Amount)
	require.Equal(t, "$19.99", subItem.NextPayment.Amount.AmountFormatted)
	require.Equal(t, int64(1672531200), subItem.NextPayment.Date)
}

func TestBillingPaymentMethod_WithCardTypeAndExpiryFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_123","last4":"4242","status":"active","card_type":"visa","expiry_year":2025,"expiry_month":12,"is_removable":true}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonData), &paymentMethod)
	require.NoError(t, err)
	require.Equal(t, "pm_123", paymentMethod.ID)
	require.Equal(t, "card", paymentMethod.PaymentType)

	// CardType, ExpiryYear, and ExpiryMonth are now nullable pointers
	require.NotNil(t, paymentMethod.CardType, "CardType should not be nil when present")
	require.Equal(t, "visa", *paymentMethod.CardType)
	require.NotNil(t, paymentMethod.ExpiryYear, "ExpiryYear should not be nil when present")
	require.Equal(t, 2025, *paymentMethod.ExpiryYear)
	require.NotNil(t, paymentMethod.ExpiryMonth, "ExpiryMonth should not be nil when present")
	require.Equal(t, 12, *paymentMethod.ExpiryMonth)

	// IsRemovable is non-nullable
	require.True(t, paymentMethod.IsRemovable)
}

func TestPlan_WithNullableDescriptionAndAvatarURL(t *testing.T) {
	t.Parallel()

	jsonDataWithOptionals := `{"object":"plan","id":"plan_123","name":"Premium Plan","description":"A premium tier plan","avatar_url":"https://example.com/avatar.png","product_id":"prod_123","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":false,"slug":"premium","free_trial_enabled":true}`
	var plan clerk.Plan
	err := json.Unmarshal([]byte(jsonDataWithOptionals), &plan)
	require.NoError(t, err)
	require.Equal(t, "plan_123", plan.ID)

	// Description and AvatarURL are now nullable pointers
	require.NotNil(t, plan.Description, "Description should not be nil when present")
	require.Equal(t, "A premium tier plan", *plan.Description)
	require.NotNil(t, plan.AvatarURL, "AvatarURL should not be nil when present")
	require.Equal(t, "https://example.com/avatar.png", *plan.AvatarURL)

	// Test Plan without description and avatar_url - use a fresh variable
	jsonDataWithoutOptionals := `{"object":"plan","id":"plan_456","name":"Basic Plan","product_id":"prod_123","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":false,"slug":"basic","free_trial_enabled":false}`
	var plan2 clerk.Plan
	err = json.Unmarshal([]byte(jsonDataWithoutOptionals), &plan2)
	require.NoError(t, err)
	require.Equal(t, "plan_456", plan2.ID)

	// Description and AvatarURL should be nil when omitted
	require.Nil(t, plan2.Description, "Description should be nil when omitted")
	require.Nil(t, plan2.AvatarURL, "AvatarURL should be nil when omitted")
}

func TestBillingPaymentMethod_PaymentTypeFieldRename(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		paymentType string
		jsonData    string
	}{
		{
			name:        "card payment type",
			paymentType: "card",
			jsonData:    `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_123","last4":"4242","status":"active","is_removable":true}`,
		},
		{
			name:        "bank_account payment type",
			paymentType: "bank_account",
			jsonData:    `{"object":"payment_method","id":"pm_456","payer_id":"payer_789","payment_type":"bank_account","gateway":"stripe","gateway_external_id":"ext_456","last4":"1234","status":"active","is_removable":true}`,
		},
		{
			name:        "digital wallet payment type",
			paymentType: "digital_wallet",
			jsonData:    `{"object":"payment_method","id":"pm_789","payer_id":"payer_012","payment_type":"digital_wallet","wallet_type":"apple_pay","gateway":"stripe","gateway_external_id":"ext_789","last4":"5678","status":"active","is_removable":false}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var paymentMethod clerk.BillingPaymentMethod
			err := json.Unmarshal([]byte(tc.jsonData), &paymentMethod)
			require.NoError(t, err)
			require.Equal(t, tc.paymentType, paymentMethod.PaymentType, "PaymentType field should be correctly mapped from payment_type JSON field")
		})
	}
}

func TestBillingPaymentMethod_IsRemovableFalseVsOmitted(t *testing.T) {
	t.Parallel()

	jsonDataWithFalse := `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_123","last4":"4242","status":"active","is_removable":false}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonDataWithFalse), &paymentMethod)
	require.NoError(t, err)
	require.False(t, paymentMethod.IsRemovable, "IsRemovable should be false when explicitly set to false")

	// Test when IsRemovable is omitted (should default to false, the zero value)
	jsonDataOmitted := `{"object":"payment_method","id":"pm_456","payer_id":"payer_789","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_456","last4":"5555","status":"active"}`
	var paymentMethod2 clerk.BillingPaymentMethod
	err = json.Unmarshal([]byte(jsonDataOmitted), &paymentMethod2)
	require.NoError(t, err)
	require.False(t, paymentMethod2.IsRemovable, "IsRemovable should default to false (zero value) when omitted")
}

func TestPlan_EmptyFeaturesArrayVsOmitted(t *testing.T) {
	t.Parallel()

	jsonDataEmptyArray := `{"object":"plan","id":"plan_123","name":"Basic Plan","product_id":"prod_123","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":false,"slug":"basic","free_trial_enabled":false,"features":[]}`
	var plan1 clerk.Plan
	err := json.Unmarshal([]byte(jsonDataEmptyArray), &plan1)
	require.NoError(t, err)
	require.NotNil(t, plan1.Features, "Features should not be nil when explicitly set as empty array")
	require.Empty(t, plan1.Features, "Features should be empty array")
	require.Len(t, plan1.Features, 0, "Features should have length 0")

	// Test with omitted features field
	jsonDataOmitted := `{"object":"plan","id":"plan_456","name":"Pro Plan","product_id":"prod_456","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":false,"slug":"pro","free_trial_enabled":false}`
	var plan2 clerk.Plan
	err = json.Unmarshal([]byte(jsonDataOmitted), &plan2)
	require.NoError(t, err)
	require.Nil(t, plan2.Features, "Features should be nil when omitted from JSON")
}

func TestBillingProduct_EmptyPlansArrayVsOmitted(t *testing.T) {
	t.Parallel()
	jsonDataEmptyArray := `{"object":"product","id":"prod_123","slug":"basic-product","currency":"USD","name":"Basic Product","is_default":false,"plans":[]}`
	var product1 clerk.BillingProduct
	err := json.Unmarshal([]byte(jsonDataEmptyArray), &product1)
	require.NoError(t, err)
	require.NotNil(t, product1.Plans, "Plans should not be nil when explicitly set as empty array")
	require.Empty(t, product1.Plans, "Plans should be empty array")
	require.Len(t, product1.Plans, 0, "Plans should have length 0")

	// Test with omitted plans field
	jsonDataOmitted := `{"object":"product","id":"prod_456","slug":"pro-product","currency":"USD","name":"Pro Product","is_default":true}`
	var product2 clerk.BillingProduct
	err = json.Unmarshal([]byte(jsonDataOmitted), &product2)
	require.NoError(t, err)
	require.Nil(t, product2.Plans, "Plans should be nil when omitted from JSON")
}

func TestSubscriptionItem_MixedOptionalFields(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","instance_id":"ins_123","status":"active","payer_id":"payer_456","plan_id":"plan_789","plan_period":"monthly","payment_method_id":"pm_123","is_free_trial":false,"period_start":1640995200,"plan":{"object":"plan","id":"plan_789","name":"Pro Plan","product_id":"prod_123","for_payer_type":"user","is_default":false,"is_recurring":true,"publicly_visible":true,"has_base_fee":true,"slug":"pro","free_trial_enabled":false},"amount":{"amount":2999,"amount_formatted":"$29.99","currency":"USD","currency_symbol":"$"},"payer":{"object":"payer","id":"payer_456","instance_id":"ins_123","user_id":"user_456"},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
			Method: http.MethodGet,
			Path:   "/v1/billing/subscription_items",
		},
	}
	client := NewClient(config)
	subscriptionItemList, err := client.ListSubscriptionItems(context.Background(), &ListSubscriptionItemsParams{})
	require.NoError(t, err)

	subItem := subscriptionItemList.Data[0]
	require.Equal(t, "sub_item_123", subItem.ID)

	// Plan is present
	require.NotNil(t, subItem.Plan, "Plan should be present")
	require.Equal(t, "plan_789", subItem.Plan.ID)
	require.Equal(t, "Pro Plan", subItem.Plan.Name)

	// PaymentMethod is omitted
	require.Nil(t, subItem.PaymentMethod, "PaymentMethod should be nil when omitted")

	// Amount is present
	require.Equal(t, int64(2999), subItem.Amount.Amount)
	require.Equal(t, "$29.99", subItem.Amount.AmountFormatted)

	// LifetimePaid is omitted (non-nullable, so should have zero values)
	require.Equal(t, int64(0), subItem.LifetimePaid.Amount, "LifetimePaid should have zero values when omitted")
	require.Equal(t, "", subItem.LifetimePaid.AmountFormatted, "LifetimePaid AmountFormatted should be empty when omitted")

	// Payer is present (non-nullable)
	require.Equal(t, "payer_456", subItem.Payer.ID)
	require.Equal(t, "user_456", *subItem.Payer.UserID)
}

func TestBillingPaymentMethod_PartialCardData(t *testing.T) {
	t.Parallel()

	jsonDataPartial := `{"object":"payment_method","id":"pm_123","payer_id":"payer_456","payment_type":"card","gateway":"stripe","gateway_external_id":"ext_123","last4":"4242","status":"active","card_type":"visa","is_removable":true}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonDataPartial), &paymentMethod)
	require.NoError(t, err)
	require.Equal(t, "pm_123", paymentMethod.ID)
	require.Equal(t, "card", paymentMethod.PaymentType)

	// CardType is present
	require.NotNil(t, paymentMethod.CardType, "CardType should be present")
	require.Equal(t, "visa", *paymentMethod.CardType)

	// Expiry fields are omitted (nullable, so should be nil)
	require.Nil(t, paymentMethod.ExpiryYear, "ExpiryYear should be nil when omitted")
	require.Nil(t, paymentMethod.ExpiryMonth, "ExpiryMonth should be nil when omitted")
	require.Nil(t, paymentMethod.WalletType, "WalletType should be nil for card payment")
}

func TestBillingPaymentMethod_DigitalWalletWithoutCardFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payment_method","id":"pm_wallet_123","payer_id":"payer_456","payment_type":"digital_wallet","gateway":"stripe","gateway_external_id":"ext_wallet_123","last4":"8888","status":"active","wallet_type":"google_pay","is_removable":true}`
	var paymentMethod clerk.BillingPaymentMethod
	err := json.Unmarshal([]byte(jsonData), &paymentMethod)
	require.NoError(t, err)
	require.Equal(t, "pm_wallet_123", paymentMethod.ID)
	require.Equal(t, "digital_wallet", paymentMethod.PaymentType)

	// WalletType is present
	require.NotNil(t, paymentMethod.WalletType, "WalletType should be present for digital wallet")
	require.Equal(t, "google_pay", *paymentMethod.WalletType)

	// Card-specific fields should be nil
	require.Nil(t, paymentMethod.CardType, "CardType should be nil for digital wallet")
	require.Nil(t, paymentMethod.ExpiryYear, "ExpiryYear should be nil for digital wallet")
	require.Nil(t, paymentMethod.ExpiryMonth, "ExpiryMonth should be nil for digital wallet")
}

func TestPayer_MixedUserFields(t *testing.T) {
	t.Parallel()

	jsonData := `{"object":"payer","id":"payer_123","instance_id":"ins_123","user_id":"user_456","first_name":"John","email":"john@example.com"}`
	var payer clerk.Payer
	err := json.Unmarshal([]byte(jsonData), &payer)
	require.NoError(t, err)
	require.Equal(t, "payer_123", payer.ID)
	require.Equal(t, "user_456", *payer.UserID)

	// FirstName and Email are present
	require.NotNil(t, payer.FirstName, "FirstName should be present")
	require.Equal(t, "John", *payer.FirstName)
	require.NotNil(t, payer.Email, "Email should be present")
	require.Equal(t, "john@example.com", *payer.Email)

	// LastName is omitted
	require.Nil(t, payer.LastName, "LastName should be nil when omitted")

	// Org fields should be nil for user payer
	require.Nil(t, payer.OrganizationID, "OrganizationID should be nil for user payer")
	require.Nil(t, payer.OrganizationName, "OrganizationName should be nil for user payer")
}
