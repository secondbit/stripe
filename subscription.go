package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Subscription struct {
	Status      string    "status"
	Object      string    "object"
	PeriodStart int       "current_period_start"
	PeriodEnd   int       "current_period_end"
	Start       int       "start"
	TrialStart  int       "trial_start"
	TrialEnd    int       "trial_end"
	Plan        *Plan     "plan"
	Customer    string    "customer"
	Error       *RawError "error"
}

func (stripe *Stripe) Subscribe(customer, plan string) (resp *Subscription, err error) {
	return stripe.RawSubscribe(customer, plan, "", -1, "", "", "", "", "", "", "", "", "", "", true)
}

func (stripe *Stripe) SubscribeWithCoupon(customer, plan, coupon string) (resp *Subscription, err error) {
	return stripe.RawSubscribe(customer, plan, coupon, -1, "", "", "", "", "", "", "", "", "", "", true)
}

func (stripe *Stripe) SubscribeWithTrial(customer, plan string, trial_end int) (resp *Subscription, err error) {
	return stripe.RawSubscribe(customer, plan, "", trial_end, "", "", "", "", "", "", "", "", "", "", true)
}

func (stripe *Stripe) SubscribeWithCard(customer, plan, number, exp_month, exp_year, cvc, name, address1, address2, zip, state, country string) (resp *Subscription, err error) {
	return stripe.RawSubscribe(customer, plan, "", -1, number, exp_month, exp_year, cvc, name, address1, address2, zip, state, country, true)
}

func (stripe *Stripe) RawSubscribe(customer, plan, coupon string, trial_end int, number, exp_month, exp_year, cvc, name, address1, address2, zip, state, country string, prorate bool) (resp *Subscription, err error) {
	values := make(url.Values)
	values.Set("plan", plan)
	if coupon != "" {
		values.Set("coupon", coupon)
	}
	if trial_end >= 0 {
		values.Set("trial_end", strconv.Itoa(trial_end))
	}
	if prorate != true {
		values.Set("prorate", "false")
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
	r, err := stripe.request("POST", "customers/"+customer+"/subscription", data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}

func (stripe *Stripe) Unsubscribe(customer string) (resp *Subscription, err error) {
	return stripe.RawUnsubscribe(customer, false)
}

func (stripe *Stripe) RawUnsubscribe(customer string, at_period_end bool) (resp *Subscription, err error) {
	values := make(url.Values)
	if at_period_end {
		values.Set("at_period_end", "true")
	}
	params := values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := stripe.request("DELETE", "customers/"+customer+"/subscription"+params, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}
