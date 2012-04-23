package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Plan struct {
	Name      string    "name"
	Object    string    "object"
	ID        string    "id"
	Interval  string    "interval"
	Currency  string    "currency"
	Amount    int       "amount"
	TrialDays int       "trial_period_days"
	Error     *RawError "error"
}

func (stripe *Stripe) CreatePlan(id, name, interval, currency string, amount int) (resp *Plan, err error) {
	return stripe.CreatePlanWithTrial(id, name, interval, currency, amount, -1)
}

func (stripe *Stripe) CreatePlanWithTrial(id, name, interval, currency string, amount, trial int) (resp *Plan, err error) {
	values := make(url.Values)
	values.Set("id", id)
	values.Set("name", name)
	values.Set("interval", interval)
	values.Set("currency", currency)
	values.Set("amount", strconv.Itoa(amount))
	if trial >= 0 {
		values.Set("trial_period_days", strconv.Itoa(trial))
	}
	params := values.Encode()
	r, err := stripe.request("POST", "plans", params)
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

func (stripe *Stripe) GetPlan(id string) (resp *Plan, err error) {
	r, err := stripe.request("GET", "plans/"+id, "")
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

func (stripe *Stripe) UpdatePlan(id, name string) (resp *Plan, err error) {
	values := make(url.Values)
	values.Set("name", name)
	params := values.Encode()
	r, err := stripe.request("POST", "plans/"+id, params)
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

func (stripe *Stripe) DeletePlan(id string) (success bool, err error) {
	r, err := stripe.request("DELETE", "plans/"+id, "")
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

func (stripe *Stripe) ListPlans() (resp []*Plan, err error) {
	return stripe.QueryPlans(-1, -1)
}

func (stripe *Stripe) QueryPlans(count, offset int) (resp []*Plan, err error) {
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
	r, err := stripe.request("GET", "plans"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int "count"
		Data  []*Plan
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
