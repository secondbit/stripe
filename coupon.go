package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Coupon represents a coupon object, according to the Stripe API.
type Coupon struct {
	ID               string    "id"
	Duration         string    "duration"           // "forever", "once", "repeating"
	DurationInMonths int       "duration_in_months" // Only useful if Duration is "repeating"
	PercentOff       int       "percent_off"
	MaxRedemptions   int       "max_redemptions"
	RedeemBy         int64     "redeem_by"
	TimesRedeemed    int       "times_redeemed"
	Object           string    "object" // Should always be "coupon"
	Error            *RawError "error"
}

// Values assigns the applicable properties of *coupon to the appropriate keys
// in *values. This makes constructing an HTTP request around a Coupon simpler.
func (coupon *Coupon) Values(values *url.Values) error {
	if coupon == nil {
		// TODO: throw an error
	}
	if coupon.PercentOff < 1 {
		// TODO: throw an error
	} else {
		values.Set("percent_off", strconv.Itoa(coupon.PercentOff))
	}
	if coupon.Duration == "" {
		// TODO: throw an error
	} else {
		values.Set("duration", coupon.Duration)
	}
	if coupon.Duration == "repeating" {
		if coupon.DurationInMonths <= 0 {
			// TODO: throw an error
		} else {
			values.Set("duration_in_months", strconv.Itoa(coupon.DurationInMonths))
		}
	}
        if coupon.MaxRedemptions > 0 {
		values.Set("max_redemptions", strconv.Itoa(coupon.MaxRedemptions))
	}
	if coupon.RedeemBy != 0 {
		values.Set("redeem_by", strconv.FormatInt(coupon.RedeemBy, 10))
	}
	if coupon.ID != "" {
		values.Set("id", coupon.ID)
	}
	return nil
}

// Discount represents the actual application of a Coupon to a particular Customer.
// It contains information about when the Discount began and will end (if the Coupon has a set duration).
type Discount struct {
	ID       string  `json:"id"`
	Object   string  `json:"object"` // Should always be "discount"
	Coupon   *Coupon `json:"coupon"`
	Customer string  `json:"customer"` // The Customer ID
	Start    int64   `json:"start"`
	End      int64   `json:"end"`
}

// CreateCoupon creates a coupon in Stripe.
func (stripe *Stripe) CreateCoupon(coupon *Coupon) (resp *Coupon, err error) {
	values := make(url.Values)
	coupon.Values(&values)
	data := values.Encode()
	r, err := stripe.request("POST", "coupons", data)
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

// GetCoupon retrieves information about the Coupon whose ID is specified by id.
func (stripe *Stripe) GetCoupon(id string) (resp *Coupon, err error) {
	r, err := stripe.request("GET", "coupons/"+id, "")
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

// DeleteCoupon deletes the Coupon whose ID is specified by id.
// Deleting a Coupon does not affect any customers who have already applied the Coupon.
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
	if err != nil {
		return false, err
	}
	if raw.Error != nil {
		// TODO: throw an error
	}
	return raw.Success, err
}

// ListCoupons queries the server for information about all your Coupons.
//
// Both the arguments are optional.
//
// Pass -1 to count to use the Stripe default (10). Count determines the number of Coupons to return. The maximum is 100.
// 
// Pass -1 to offset to use the Stripe default (0). Offset determines the number of recent Coupons to skip.
func (stripe *Stripe) ListCoupons(count, offset int) (resp []*Coupon, err error) {
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
