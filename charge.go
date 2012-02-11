package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Charge struct {
	Amount   int    "amount"
	Currency string "currency"
	Card     struct {
		Type         string "type"
		ExpYear      int    "exp_year"
		Country      string "country"
		LastFour     string "last4"
		Object       string "object"
		ExpMonth     int    "exp_month"
		CVCCheck     string "cvc_check"
		AddressCheck string "address_line1_check"
		ZipCheck     string "address_zip_check"
	}
	Customer    string
	Description string    "description"
	Created     int       "created"
	Fee         int       "fee"
	ID          string    "id"
	LiveMode    bool      "livemode"
	Object      string    "object"
	Paid        bool      "paid"
	Refunded    bool      "refunded"
	Error       *RawError "error"
}

func (stripe *Stripe) CreateCharge(amount int, currency string) (resp *Charge, err error) {
	return stripe.RawCreateCharge(amount, currency, "", "", "", "", "", "", "", "", "", "", "", "")
}

func (stripe *Stripe) ChargeCustomer(amount int, currency, customer string) (resp *Charge, err error) {
	return stripe.RawCreateCharge(amount, currency, customer, "", "", "", "", "", "", "", "", "", "", "")
}

func (stripe *Stripe) ChargeCustomerWithDescription(amount int, currency, customer, description string) (resp *Charge, err error) {
	return stripe.RawCreateCharge(amount, currency, customer, description, "", "", "", "", "", "", "", "", "", "")
}

func (stripe *Stripe) RawCreateCharge(amount int, currency, customer, description, number, exp_month, exp_year, cvc, name, address1, address2, zip, state, country string) (resp *Charge, err error) {
	values := make(url.Values)
	values.Set("amount", strconv.Itoa(amount))
	values.Set("currency", currency)
	if customer != "" {
		values.Set("customer", customer)
	}
	if description != "" {
		values.Set("description", description)
	}
	if number != "" {
		values.Set("card[number]", number)
	}
	if exp_month != "" {
		values.Set("card[exp_month]", exp_month)
	}
	if exp_year != "" {
		values.Set("card[exp_year]", exp_year)
	}
	if cvc != "" {
		values.Set("card[cvc]", cvc)
	}
	if name != "" {
		values.Set("card[name]", name)
	}
	if address1 != "" {
		values.Set("card[address_line1]", address1)
	}
	if address2 != "" {
		values.Set("card[address_line2]", address2)
	}
	if zip != "" {
		values.Set("card[address_zip]", zip)
	}
	if state != "" {
		values.Set("card[address_state]", state)
	}
	if country != "" {
		values.Set("card[address_country]", country)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "charges", data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
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
		Count int "count"
		Data  []*Charge
		Error *RawError "error"
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
