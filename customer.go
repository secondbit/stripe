package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// Customer represents a customer, according to the Stripe API
type Customer struct {
	Description  string        `json:"description"`
	Object       string        `json:"object"` // Should always be "customer"
	LiveMode     bool          `json:"livemode"`
	ActiveCard   *Card         `json:"active_card"`
	Created      int64         `json:"created"`
	ID           string        `json:"id"`
	Balance      int64         `json:"account_balance"`
	Error        *RawError     `json:"error"`
	Delinquent   bool          `json:"delinquent"` // Whether the latest charge failed
	Discount     *Discount     `json:"discount"`
	Email        string        `json:"email"`
	Subscription *Subscription `json:"subscription"`
}

// ChargeValues sets *customer's non-empty properties to their appropriate key in *values
// This is useful for constructing HTTP requests from Customer objects
// This also satisfies the Chargeable interface, allowing Customers to be charged
func (customer *Customer) ChargeValues(values *url.Values) error {
	if customer == nil {
		// TODO: Throw an error
	}
	if customer.ActiveCard == nil {
		// TODO: Throw an error
	}
	if customer.ID != "" {
		values.Set("customer", customer.ID)
	} else {
		// TODO: throw an error
	}
	return nil
}

func (customer *Customer) Values(values *url.Values) error {
	if customer == nil {
		// TODO: Throw an error
	}
	if customer.Description != "" {
		values.Set("description", customer.Description)
	}
	if customer.Email != "" {
		values.Set("email", customer.Email)
	}
	return nil
}

// CreateCustomer creates a new customer on Stripe.
//
// All arguments except the Customer object are optional.
//
// If a Chargeable is provided, it is associated with that customer and automatically validated. Pass nil to omit the Chargeable.
//
// If plan is non-empty, it is used as the identifier of a Plan to subscribe customer to.
//
// If coupon is non-empty, it is used as a Coupon that will be applied to all of customer's recurring charges.
//
// If trial_end is not -1, it will be used as the UTC integer timestamp representing thend of the trial period for customer.
//
// trial_end overrides the plan's default trial period, if not -1.
func (stripe *Stripe) CreateCustomer(customer *Customer, chargeable Chargeable, plan, coupon string, trial_end int64) (resp *Customer, err error) {
	values := make(url.Values)
	err = customer.Values(&values)
	if err != nil {
		return nil, err
	}
	if chargeable != nil {
		err = chargeable.ChargeValues(values)
		if err != nil {
			return nil, err
		}
	}
	if plan != "" {
		values.Set("plan", plan)
	}
	if trial_end != -1 {
		values.Set("trial_end", strconv.FormatInt(trial_end, 10))
	}
	if coupon != "" {
		values.Set("coupon", coupon)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "customers", data)
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

// UpdateCustomer updates the information Stripe has on *customer.
//
// All arguments except the Customer object are optional.
//
// If a Chargeable is provided, it is attached to *customer and automatically validated. Pass nil to omit the Chargeable.
//
// If couponID is non-empty, it is used as a Coupon that will be applied to all of *customer's recurring charges.
func (stripe *Stripe) UpdateCustomer(customer *Customer, chargeable Chargeable, couponID string) (resp *Customer, err error) {
	if customer.ID == "" {
		// TODO: throw an error
	}
	values := make(url.Values)
	err = customer.Values(&values)
	if err != nil {
		return nil, err
	}
	if chargeable != nil {
		err = chargeable.ChargeValues(&values)
		if err != nil {
			return nil, err
		}
	}
	if couponID != "" {
		values.Set("coupon", couponID)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "customers/"+customer.ID, data)
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

// GetCustomer retrieves information on the Customer with ID of id.
func (stripe *Stripe) GetCustomer(id string) (resp *Customer, err error) {
	r, err := stripe.request("GET", "customers/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
	}
	return resp, err
}

// DeleteCustomer permanently deletes the Customer with ID of id from Stripe. It cannot be undone.
func (stripe *Stripe) DeleteCustomer(id string) (success bool, err error) {
	r, err := stripe.request("DELETE", "customers/"+id, "")
	if err != nil {
		return false, err
	}
	var raw struct {
		Success bool      "deleted"
		ID      string    "id"
		Error   *RawError "error"
	}
	err = json.Unmarshal(r, &raw)
	if err != nil {
		return false, err
	}
	if raw.Error != nil {
		// TODO: throw an error
	}
	return raw.Success, err
}

// ListCustomers queries the server for information about all your Customers. Results are returned sorted by creation date, with the most recently created Customers appearing first.
//
// Both the arguments are optional.
//
// Pass -1 to count to use the Stripe default (10). Count determines the number of customers to return. The maximum is 100.
// 
// Pass -1 to offset to use the Stripe default (0). Offset determines the number of recent customers to skip.
func (stripe *Stripe) ListCustomers(count, offset int) (resp []*Customer, err error) {
	values := make(url.Values)
	if count >= 0 {
		values.Set("count", strconv.Itoa(count))
	}
	if offset >= 0 {
		values.Set("offset", strconv.Itoa(offset))
	}
	params := values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := stripe.request("GET", "customers"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int         "count"
		Data  []*Customer "data"
		Error *RawError   "error"
	}
	err = json.Unmarshal(r, &raw)
	if err != nil {
		return nil, err
	}
	if raw.Error != nil {
		// TODO: throw an error
	}
	resp = raw.Data
	return
}
