package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/clerk/clerk-sdk-go/v3"
	"github.com/clerk/clerk-sdk-go/v3/clerktest"
	"github.com/stretchr/testify/require"
)

func TestBillingClientListPlans(t *testing.T) {
	t.Parallel()
	config := &clerk.ClientConfig{}
	config.HTTPClient = &http.Client{
		Transport: &clerktest.RoundTripper{
			T:      t,
			Out:    json.RawMessage(`{"data": [{"object":"plan","id":"plan_123","name":"Basic Plan","for_payer_type":"user","features":[{"object":"feature","id":"feature_456","name":"Feature 1","key":"feature_1"}],"unit_prices":[{"name":"Seats","block_size":1,"tiers":[{"starts_at_block":1,"ends_after_block":10,"fee_per_block":{"amount":500,"amount_formatted":"$5.00","currency":"usd","currency_symbol":"$"}},{"starts_at_block":11,"ends_after_block":null,"fee_per_block":{"amount":400,"amount_formatted":"$4.00","currency":"usd","currency_symbol":"$"}}]}]}],"total_count": 1}`),
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
	require.Equal(t, 1, len(planList.Data[0].UnitPrices))
	require.Equal(t, "Seats", planList.Data[0].UnitPrices[0].Name)
	require.Equal(t, int64(1), planList.Data[0].UnitPrices[0].BlockSize)
	require.Equal(t, 2, len(planList.Data[0].UnitPrices[0].Tiers))
	require.Equal(t, int64(1), planList.Data[0].UnitPrices[0].Tiers[0].StartsAtBlock)
	require.Equal(t, int64(10), *planList.Data[0].UnitPrices[0].Tiers[0].EndsAfterBlock)
	require.Equal(t, int64(500), planList.Data[0].UnitPrices[0].Tiers[0].FeePerBlock.Amount)
	require.Equal(t, "$5.00", planList.Data[0].UnitPrices[0].Tiers[0].FeePerBlock.AmountFormatted)
	require.Equal(t, int64(11), planList.Data[0].UnitPrices[0].Tiers[1].StartsAtBlock)
	require.Nil(t, planList.Data[0].UnitPrices[0].Tiers[1].EndsAfterBlock)
	require.Equal(t, int64(400), planList.Data[0].UnitPrices[0].Tiers[1].FeePerBlock.Amount)
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
			Out:    json.RawMessage(`{"data": [{"object":"subscription_item","id":"sub_item_123","payer_id":"payer_456","plan_id":"plan_789","status":"active","period_start":1640995200,"period_end":1643673600,"payer":{"object":"payer","id":"payer_456","user_id":"user_456","first_name":"John","last_name":"Doe","email":"john@example.com","created_at":1640995200,"updated_at":1640995200},"plan":{"object":"plan","id":"plan_789","name":"Pro Plan","payer_type":["user"],"features":[]},"seats":{"quantity":10},"created_at":1640995200,"updated_at":1640995200}],"total_count": 1}`),
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
	require.NotNil(t, subscriptionItemList.Data[0].Seats)
	require.Equal(t, int64(10), *subscriptionItemList.Data[0].Seats.Quantity)
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
