package stripe

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// Chargeable is an interface to expose items that can be charged.
// The Stripe API uses credit cards, tokens, and customers interchangeably for most operations
// dealing with money; rather than writing three versions of all these operations, Chargeable 
// allows a single version to work for any of the types that satisfy it.
type Chargeable interface {
	ChargeValues(values *url.Values) error
}

type Charge struct {
	Amount      int       `json:"amount"`
	Currency    string    `json:"currency"`
	Card        *Card     `json:"card"`
	Customer    string    `json:"customer"` // The Customer's ID
	Description string    `json:"description"`
	Created     int       `json:"created"`
	Fee         int       `json:"fee"`
	ID          string    `json:"id"`
	LiveMode    bool      `json:"livemode"`
	Object      string    `json:"object"` // Should always be "charge"
	Paid        bool      `json:"paid"`
	Refunded    bool      `json:"refunded"`
	Error       *RawError `json:"error"`
}

// CreateCharge submits a charge object to the Stripe servers, at which point Stripe will charge the card.
// Description is optional.
func (stripe *Stripe) CreateCharge(chargeable Chargeable, amount int, currency, description string) (resp *Charge, err error) {
	values := make(url.Values)
	values.Set("amount", strconv.Itoa(amount))
	values.Set("currency", currency)
	if description != "" {
		values.Set("description", description)
	}
	err = chargeable.ChargeValues(&values)
	if err != nil {
		return nil, err
	}
	data := values.Encode()
	r, err := stripe.request("POST", "charges", data)
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

// GetCharge retrieves the details of a charge that has previously been created.
func (stripe *Stripe) GetCharge(id string) (resp *Charge, err error) {
	if id == "" {
		return nil, errors.New("No ID set.")
	}
	r, err := stripe.request("GET", "charges/"+id, "")
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
	return resp, err
}

// RefundCharge refunds all or part of a charge. To refund all of a charge, pass -1
// as the amount.
func (stripe *Stripe) RefundCharge(id string, amount int) (resp *Charge, err error) {
	var body string
	if amount >= 0 {
		values := make(url.Values)
		values.Set("amount", strconv.Itoa(amount))
		body = values.Encode()
	}
	r, err := stripe.request("POST", "charges/"+id+"/refund", body)
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

// ListCharges queries the server for information about past charges.
//
// All the arguments are optional.
//
// Pass -1 to count to use the Stripe default (10). Count determines the number of charges to return. The maximum is 100.
// 
// Pass -1 to offset to use the Stripe default (0). Offset determines the number of recent charges to skip.
// 
// Pass anything but an empty string to customer to show only that customer's charges.
//
func (stripe *Stripe) ListCharges(count, offset int, customer string) (resp []*Charge, err error) {
	values := make(url.Values)
	if count >= 0 {
		values.Set("count", strconv.Itoa(count))
	}
	if offset >= 0 {
		values.Set("offset", strconv.Itoa(offset))
	}
	if customer != "" {
		values.Set("customer", customer)
	}
	params := values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := stripe.request("GET", "charges"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int "count"
		Data  []*Charge
		Error *RawError "error"
	}
	err = json.Unmarshal(r, &raw)
	if err != nil {
		return nil, err
	}
	if raw.Error != nil {
		//TODO: Throw an error
	}
	resp = raw.Data
	return
}
