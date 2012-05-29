package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Invoice struct {
        ID                 string  `json:"id"`
        LiveMode           bool    `json:"livemode"`
        AmountDue          int     `json:"amount_due"`
        AttemptCount       int     `json:"attempt_count"`
        Attempted          bool    `json:"attempted"`
        Closed             bool    `json:"closed"`
        CustomerID         string  `json:"customer"`
        Date               int64   `json:"date"`
        Paid               bool    `json:"paid"`
        PeriodEnd          int64   `json:"period_end"`
        PeriodStart        int64   `json:"period_start"`
        StartingBalance    int     `json:"starting_balance"`
        Subtotal           int     `json:"subtotal"`
        Total              int     `json:"total"`
        ChargeID           *string `json:"charge"`
        Discount           *Discount `json:"discount"`
        EndingBalance      *int    `json:"ending_balance"`
        NextPaymentAttempt *int    `json:"next_payment_attempt"`
	Lines              struct {
                InvoiceItems  []*InvoiceItem `json:"invoiceitems"`
                Subscriptions []*SubscriptionItem `json:"subscriptions"`
                Prorated      []*InvoiceItem `json:"prorations"`
	}
        Object string    `json:"object"`
        Error  *RawError `json:"error"`
}

type SubscriptionItem struct {
        Amount             int     `json:"amount"`
        Period             struct {
                Start      int64   `json:"start"`
                End        int64   `json:"end"`
        }                          `json:"period"`
        Plan               *Plan   `json:"plan"`
}

func (stripe *Stripe) GetInvoice(id string) (resp *Invoice, err error) {
	r, err := stripe.request("GET", "invoices/"+id, "")
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

func (stripe *Stripe) GetNextInvoice(customer string) (resp *Invoice, err error) {
	values := make(url.Values)
	values.Set("customer", customer)
	params := values.Encode()
	r, err := stripe.request("GET", "invoices/upcoming"+params, "")
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

func (stripe *Stripe) ListInvoices(count, offset int, customer string) (resp []*Invoice, err error) {
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
	r, err := stripe.request("GET", "invoices"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int        "count"
		Data  []*Invoice "data"
		Error *RawError  "error"
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

type InvoiceItem struct {
        ID          string    `json:"id"`
        LiveMode    bool      `json:"livemode"`
        Date        int64     `json:"date"`
        Description *string   `json:"description"`
        Currency    string    `json:"currency"`
        Amount      int       `json:"amount"`
        CustomerID  string    `json:"customer"`
        InvoiceID   *string   `json:"invoice"`
        Object      int       `json:"object"`
        Error       *RawError `json:"error"`
}

func (item *InvoiceItem) Values(values *url.Values) error {
        if item == nil {
                //TODO: Throw error
        }
        if item.CustomerID == "" {
                // TODO: Throw error
        }
        if item.Amount == 0 {
                // TODO: Throw error
        }
        if item.Currency != "USD" {
                // TODO: Throw error
        }
        if item.InvoiceID != nil {
                values.Set("invoice", *item.InvoiceID)
        }
        if item.Description != nil {
                values.Set("desription", *item.Description)
        }
        values.Set("customer", item.CustomerID)
        values.Set("amount", strconv.Itoa(item.Amount))
        values.Set("currency", item.Currency)
        return nil
}

func (stripe *Stripe) CreateInvoiceItem(item *InvoiceItem) (resp *InvoiceItem, err error) {
        values := make(url.Values)
        err = item.Values(&values)
        if err != nil {
                return nil, err
        }
	data := values.Encode()
	r, err := stripe.request("POST", "invoiceitems", data)
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

func (stripe *Stripe) GetInvoiceItem(id string) (resp *InvoiceItem, err error) {
	r, err := stripe.request("GET", "invoiceitems/"+id, "")
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

func (stripe *Stripe) UpdateInvoiceItem(id string, amount int, description string) (resp *InvoiceItem, err error) {
	values := make(url.Values)
	if amount >= 0 {
		values.Set("amount", strconv.Itoa(amount))
	}
	if description != "" {
		values.Set("description", description)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "invoiceitems/"+id, data)
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

func (stripe *Stripe) DeleteInvoiceItem(id string) (success bool, err error) {
	r, err := stripe.request("DELETE", "invoiceitems/"+id, "")
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

func (stripe *Stripe) ListInvoiceItems(count, offset int, customer string) (resp []*InvoiceItem, err error) {
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
	r, err := stripe.request("GET", "invoiceitems"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
                Count int            `json:"count"`
                Data  []*InvoiceItem `json:"data"`
                Error *RawError      `json:"error"`
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
