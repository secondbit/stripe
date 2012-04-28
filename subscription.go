package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Subscription struct {
	Status            string    `json:"status"` // "trialing"/"active"/"past_due"/"canceled"/"unpaid"
	Object            string    `json:"object"` // Should always be "subscription"
	PeriodStart       int64     `json:"current_period_start"`
	PeriodEnd         int64     `json:"current_period_end"`
	CancelAtPeriodEnd bool      `json:"cancel_at_period_end"`
	CanceledAt        int64     `json:"canceled_at"`
	EndedAt           int64     `json:"ended_at"`
	Start             int64     `json:"start"`
	TrialStart        int64     `json:"trial_start"`
	TrialEnd          int64     `json:"trial_end"`
	Plan              *Plan     `json:"plan"`
	Customer          string    `json:"customer"` // The Customer's ID
	Error             *RawError `json:"error"`
}

// Values assigns the applicable properties of *subscription to the appropriate keys
// in *values. This makes constructing an HTTP request around a Subscription simpler.
func (subscription *Subscription) Values(values *url.Values) error {
	if subscription == nil {
		// TODO: throw an error
	}
	if subscription.Plan == nil {
		// TODO: throw an error
	}
	if subscription.Plan.ID == "" {
		// TODO: throw an error
	}
	values.Set("plan", subscription.Plan.ID)
	if subscription.TrialEnd > 0 {
		values.Set("trial_end", strconv.FormatInt(subscription.TrialEnd, 10))
	}
}

// Subscribe updates the customer's plan. The customer will be billed monthly according to the new plan.
//
// *subscription is the only required argument. *subscription.Plan.ID and *subscription.Customer.ID must be set.
//
// If couponID is non-empty, it will be used as the ID of a coupon to apply to the customer.
//
// If prorate is true, the customer will be prorated to make up for the price changes
//
// If chargeable is non-nil, it will be attached to the customer. Can be either a token or a credit card.
func (stripe *Stripe) Subscribe(subscription *Subscription, couponID string, prorate bool, chargeable Chargeable) (resp *Subscription, err error) {
	values := make(url.Values)
	if subscription.Customer == nil || subscription.Customer.ID == "" {
		// TODO: throw an error
	}
	err = subscription.Values(&values)
	if err != nil {
		return nil, err
	}
	if couponID != "" {
		values.Set("coupon", couponID)
	}
	if prorate != true {
		values.Set("prorate", "false")
	}
	if chargeable != nil {
		chargeable.ChargeValues(&values)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "customers/"+subscription.Customer.ID+"/subscription", data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		// TODO: Throw an error
	}
	return
}

// Unsubscribe cancels the subscription of the customer whose ID matches customerID.
//
// If at_period_end is true, the subscription will remain active until the end of the period, at which point
// it will be cancelled and not renewed. Otherwise, the subscription will be cancelled immediately.
//
// Any pending invoice items will still be charged for at the end of the period unless they are manually deleted.
func (stripe *Stripe) Unsubscribe(customerID string, at_period_end bool) (resp *Subscription, err error) {
	values := make(url.Values)
	if at_period_end {
		values.Set("at_period_end", "true")
	}
	params := values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := stripe.request("DELETE", "customers/"+customerID+"/subscription"+params, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		// TODO: Throw an error
	}
	return
}
