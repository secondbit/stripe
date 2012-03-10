package stripe

import (
	"encoding/json"
	"net/url"
	"fmt"
        "strconv"
        "time"
)

type Customer struct {
	Description string "description"
	Object      string "object"
	LiveMode    bool   "livemode"
	ActiveCard  struct {
		ExpYear  int    "exp_year"
		Type     string "type"
		Country  string "country"
		LastFour string "last4"
		Object   string "object"
		ExpMonth string "exp_month"
	}
	Created int       "created"
	ID      string    "id"
	Error   *RawError "error"
}

func (stripe *Stripe) CreateCustomerWithToken(email, token, description, plan string, trial_end time.Time, coupon string) (resp *Customer, err error) {
        values := make(url.Values)
        if email != "" {
                values.Set("email", email)
        }
        if token != "" {
                values.Set("card", token)
        }
        if description != "" {
                values.Set("description", description)
        }
        if plan != "" {
                values.Set("plan", plan)
        }
        if trial_end.IsZero() {
                // TODO: There has to be a better way to convert int64 to string
                values.Set("trial_end", fmt.Sprintf("%v", trial_end.Unix()))
        }
        if coupon != "" {
                values.Set("coupon", "")
        }
        data := values.Encode()
        r, err := stripe.request("POST", "customers", data)
        if err != nil {
                return nil, err
        }
        err = json.Unmarshal(r, &resp)
        return
}

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
	if raw.Error != nil {
		// TODO: throw an error
	}
	return raw.Success, err
}

func (stripe *Stripe) ListCustomers() (resp []*Customer, err error) {
	return stripe.QueryCustomers(-1, -1)
}

func (stripe *Stripe) QueryCustomers(count, offset int) (resp []*Customer, err error) {
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
