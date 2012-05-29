package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Plan struct {
        Name      string    `json:"name"`
        Object    string    `json:"object"`
        ID        string    `json:"id"`
        Interval  string    `json:"interval"`
        Currency  string    `json:"currency"`
        Amount    int       `json:"amount"`
        TrialDays int       `json:"trial_period_days"`
        Error     *RawError `json:"error"`
}

func (plan *Plan) Values(values *url.Values) error {
        if plan == nil {
                // TODO: Throw error
        }
        if plan.ID == "" {
                // TODO: Throw error
        }
        if plan.Name == "" {
                // TODO: Throw error
        }
        if plan.Amount <= 0 {
                // TODO: Throw error
        }
        if plan.Currency != "USD" {
                // TODO: Throw error
        }
        if plan.Interval != "month" && plan.Interval != "year" {
                // TODO: Throw error
        }
        if plan.TrialDays > 0 {
                values.Set("trial_period_days", strconv.Itoa(plan.TrialDays))
        }
        values.Set("id", plan.ID)
        values.Set("amount", strconv.Itoa(plan.Amount))
        values.Set("currency", plan.Currency)
        values.Set("interval", plan.Interval)
        values.Set("name", plan.Name)
        return nil
}

func (stripe *Stripe) CreatePlan(plan *Plan) (resp *Plan, err error) {
	values := make(url.Values)
        if plan == nil {
                // TODO: Throw error
        }
        err = plan.Values(&values)
        if err != nil {
                return nil, err
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

func (stripe *Stripe) ListPlans(count, offset int) (resp []*Plan, err error) {
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
