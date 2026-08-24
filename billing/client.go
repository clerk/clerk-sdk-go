// Package billing provides the Billing API.
package billing

import (
	"context"
	"net/http"
	"net/url"

	"github.com/clerk/clerk-sdk-go/v2"
)

//go:generate go run ../cmd/gen/main.go

const path = "/billing"

// Client is used to invoke the Billing API.
type Client struct {
	Backend clerk.Backend
}

func NewClient(config *clerk.ClientConfig) *Client {
	return &Client{
		Backend: clerk.NewBackend(&config.BackendConfig),
	}
}

type ListPlansParams struct {
	clerk.APIParams
	clerk.ListParams
	PayerType *string
}

func (params *ListPlansParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.PayerType != nil {
		q.Set("payer_type", *params.PayerType)
	}
	return q
}

// ListPlans returns a list of billing plans.
func (c *Client) ListPlans(ctx context.Context, params *ListPlansParams) (*clerk.PlanList, error) {
	plansPath, err := clerk.JoinPath(path, "/plans")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, plansPath)
	req.SetParams(params)
	planList := &clerk.PlanList{}
	err = c.Backend.Call(ctx, req, planList)
	return planList, err
}

type ListSubscriptionItemsParams struct {
	clerk.APIParams
	clerk.ListParams
	Status         *string
	PayerType      *string
	PlanID         *string
	IncludeFree    *bool
	Query          *string
	UserID         *string
	OrganizationID *string
}

func (params *ListSubscriptionItemsParams) ToQuery() url.Values {
	q := params.ListParams.ToQuery()
	if params.Status != nil {
		q.Set("status", *params.Status)
	}
	if params.PayerType != nil {
		q.Set("payer_type", *params.PayerType)
	}
	if params.PlanID != nil {
		q.Set("plan_id", *params.PlanID)
	}
	if params.IncludeFree != nil {
		if *params.IncludeFree {
			q.Set("include_free", "true")
		} else {
			q.Set("include_free", "false")
		}
	}
	if params.Query != nil {
		q.Set("query", *params.Query)
	}
	if params.UserID != nil {
		q.Set("user_id", *params.UserID)
	}
	if params.OrganizationID != nil {
		q.Set("organization_id", *params.OrganizationID)
	}
	return q
}

// ListSubscriptionItems returns a list of subscription items.
func (c *Client) ListSubscriptionItems(ctx context.Context, params *ListSubscriptionItemsParams) (*clerk.SubscriptionItemList, error) {
	subscriptionItemsPath, err := clerk.JoinPath(path, "/subscription_items")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodGet, subscriptionItemsPath)
	req.SetParams(params)
	subscriptionItemList := &clerk.SubscriptionItemList{}
	err = c.Backend.Call(ctx, req, subscriptionItemList)
	return subscriptionItemList, err
}

type CancelSubscriptionItemParams struct {
	clerk.APIParams
	EndNow *bool
}

func (params *CancelSubscriptionItemParams) ToQuery() url.Values {
	q := url.Values{}
	if params.EndNow != nil {
		if *params.EndNow {
			q.Set("end_now", "true")
		} else {
			q.Set("end_now", "false")
		}
	}
	return q
}

// CancelSubscriptionItem cancels a subscription item.
func (c *Client) CancelSubscriptionItem(ctx context.Context, subscriptionItemID string, params *CancelSubscriptionItemParams) (*clerk.SubscriptionItem, error) {
	cancelPath, err := clerk.JoinPath(path, "/subscription_items", subscriptionItemID)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := params.ToQuery()
		if len(query) > 0 {
			cancelPath += "?" + query.Encode()
		}
	}

	req := clerk.NewAPIRequest(http.MethodDelete, cancelPath)
	subscriptionItem := &clerk.SubscriptionItem{}
	err = c.Backend.Call(ctx, req, subscriptionItem)
	return subscriptionItem, err
}

type ExtendFreeTrialParams struct {
	clerk.APIParams
	ExtendTo string `json:"extend_to"`
}

// ExtendFreeTrial extends the free trial period of a subscription item.
func (c *Client) ExtendFreeTrial(ctx context.Context, subscriptionItemID string, params *ExtendFreeTrialParams) (*clerk.SubscriptionItem, error) {
	extendPath, err := clerk.JoinPath(path, "/subscription_items", subscriptionItemID, "/extend_free_trial")
	if err != nil {
		return nil, err
	}
	req := clerk.NewAPIRequest(http.MethodPost, extendPath)
	req.SetParams(params)
	subscriptionItem := &clerk.SubscriptionItem{}
	err = c.Backend.Call(ctx, req, subscriptionItem)
	return subscriptionItem, err
}
