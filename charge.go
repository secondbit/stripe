package stripe

import (
        "encoding/json"
        "net/url"
        "strconv"
)

type Charge struct {
	Amount      int                 "amount"
	Currency    string              "currency"
	Card        struct {
                Type            string  "type"
                ExpYear         int     "exp_year"
                CVCCheck        string  "cvc_check"
                Country         string  "country"
                LastFour        string  "last4"
                Object          string  "object"
                ExpMonth        int     "exp_month"
        }
	Customer    string 
	Description string              "description"
	Created     int                 "created"
	Fee         int                 "fee"
	ID          string              "id"
	LiveMode    bool                "livemode"
	Object      string              "object"
	Paid        bool                "paid"
	Refunded    bool                "refunded"
        Error       *RawError           "error"
}

func (stripe *Stripe) GetCharge(id string) (resp *Charge, err error) {
	if id == "" {
		// TODO: throw an error
	}
	r, err := stripe.request("GET", "charges/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return resp, err
}

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
	return
}

func (stripe *Stripe) ListCharges() (resp []*Charge, err error) {
        return stripe.QueryCharges(-1, -1, "")
}

func (stripe *Stripe) ListChargesByCustomer(customer string) (resp []*Charge, err error) {
        return stripe.QueryCharges(-1, -1, customer)
}

func (stripe *Stripe) QueryCharges(count, offset int, customer string) (resp []*Charge, err error) {
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
                Count   int    "count"
                Data    []*Charge
        }
        err = json.Unmarshal(r, &raw)
        if err != nil {
                return nil, err
        }
        resp = raw.Data
	return
}
