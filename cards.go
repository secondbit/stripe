package stripe

import (
	"encoding/json"
	"net/url"
)

type CardToken struct {
	Used     bool   "used"
	Currency string "currency"
	Object   string "object"
	LiveMode bool   "livemode"
	Card     struct {
		Type     string "type"
		ExpYear  int    "exp_year"
		CVCCheck string "cvc_check"
		Country  string "country"
		LastFour string "last4"
		Object   string "object"
		ExpMonth int    "exp_month"
	}
	Created int       "created"
	ID      string    "id"
	Error   *RawError "error"
}

func (stripe *Stripe) GetCardToken(id string) (resp *CardToken, err error) {
	r, err := stripe.request("GET", "tokens/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return resp, err
}

func (stripe *Stripe) CreateCardToken(number, exp_month, exp_year string) (resp *CardToken, err error) {
	return stripe.CreateCardTokenWithAll(number, exp_month, exp_year, "", "", "", "", "", "", "")
}

func (stripe *Stripe) CreateCardTokenWithCVC(number, exp_month, exp_year, cvc string) (resp *CardToken, err error) {
	return stripe.CreateCardTokenWithAll(number, exp_month, exp_year, cvc, "", "", "", "", "", "")
}

func (stripe *Stripe) CreateCardTokenWithAll(number, exp_month, exp_year, cvc, name, address1, address2, zip, state, country string) (resp *CardToken, err error) {
	values := make(url.Values)
	values.Set("card[number]", number)
	values.Set("card[exp_month]", exp_month)
	values.Set("card[exp_year]", exp_year)
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
	params := values.Encode()
	r, err := stripe.request("POST", "tokens", params)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}
