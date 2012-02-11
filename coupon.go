package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Coupon struct {
	ID               string    "id"
	Duration         string    "duration"
	DurationInMonths int       "duration_in_months"
	PercentOff       int       "percent_off"
	Object           string    "object"
	Error            *RawError "error"
}

func (stripe *Stripe) CreateCoupon(id, duration string, percent_off int) (resp *Coupon, err error) {
	return stripe.RawCreateCoupon(id, duration, percent_off, -1, -1, "")
}

func (stripe *Stripe) RawCreateCoupon(id, duration string, percent_off, duration_in_months, max_redemptions int, redeem_by string) (resp *Coupon, err error) {
	values := make(url.Values)
	values.Set("duration", duration)
	values.Set("percent_off", strconv.Itoa(percent_off))
	if id != "" {
		values.Set("id", id)
	}
	if duration_in_months >= 0 {
		values.Set("duration_in_months", strconv.Itoa(duration_in_months))
	}
	if max_redemptions >= 0 {
		values.Set("max_redemptions", strconv.Itoa(max_redemptions))
	}
	if redeem_by != "" {
		values.Set("redeem_by", redeem_by)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "coupons", data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}

func (stripe *Stripe) GetCoupon(id string) (resp *Coupon, err error) {
	r, err := stripe.request("GET", "coupons/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}

func (stripe *Stripe) DeleteCoupon(id string) (success bool, err error) {
	r, err := stripe.request("DELETE", "coupons/"+id, "")
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

func (stripe *Stripe) ListCoupons() (resp []*Coupon, err error) {
	return stripe.QueryCoupons(-1, -1)
}

func (stripe *Stripe) QueryCoupons(count, offset int) (resp []*Coupon, err error) {
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
	r, err := stripe.request("GET", "coupons"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int       "count"
		Data  []*Coupon "data"
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
