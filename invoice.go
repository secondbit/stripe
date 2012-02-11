package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Invoice struct {
	ID                 string "id"
	Subtotal           int    "subtotal"
	Total              int    "total"
	Created            int    "created"
	NextPaymentAttempt string "next_payment_attempt"
	Lines              struct {
		InvoiceItems  []*InvoiceItem "invoiceitems"
		Subscriptions []*InvoiceLine "subscriptions" // TODO: Is this right?
		Prorated      []*InvoiceItem "prorated"      // TODO: Is this right?
	}
	Object string    "object"
	Error  *RawError "error"
}

type InvoiceLine struct {
	ID       string "id"
	Date     int    "date"
	Amount   int    "amount"
	Currency string "currency"
	Object   string "object"
}

func (stripe *Stripe) GetInvoice(id string) (resp *Invoice, err error) {
	r, err := stripe.request("GET", "invoices/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
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
	return
}

func (stripe *Stripe) ListInvoices() (resp []*Invoice, err error) {
	return stripe.QueryInvoices(-1, -1, "")
}

func (stripe *Stripe) ListInvoicesByCustomer(customer string) (resp []*Invoice, err error) {
	return stripe.QueryInvoices(-1, -1, customer)
}

func (stripe *Stripe) QueryInvoices(count, offset int, customer string) (resp []*Invoice, err error) {
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
	ID          string    "id"
	Date        int       "date"
	Description string    "description"
	Currency    string    "currency"
	Amount      int       "amount"
	Object      int       "object"
	Error       *RawError "error"
}

func (stripe *Stripe) CreateInvoiceItem(customer string, amount int, currency string) (resp *InvoiceItem, err error) {
	return stripe.RawCreateInvoiceItem(customer, "", amount, currency, "")
}

func (stripe *Stripe) CreateInvoiceItemOnInvoice(customer, invoice string, amount int, currency string) (resp *InvoiceItem, err error) {
	return stripe.RawCreateInvoiceItem(customer, invoice, amount, currency, "")
}

func (stripe *Stripe) CreateInvoiceItemWithDescription(customer string, amount int, currency, description string) (resp *InvoiceItem, err error) {
	return stripe.RawCreateInvoiceItem(customer, "", amount, currency, description)
}

func (stripe *Stripe) RawCreateInvoiceItem(customer, invoice string, amount int, currency, description string) (resp *InvoiceItem, err error) {
	values := make(url.Values)
	values.Set("customer", customer)
	values.Set("amount", strconv.Itoa(amount))
	values.Set("currency", currency)
	if invoice != "" {
		values.Set("invoice", invoice)
	}
	if description != "" {
		values.Set("description", description)
	}
	data := values.Encode()
	r, err := stripe.request("POST", "invoiceitems", data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}

func (stripe *Stripe) GetInvoiceItem(id string) (resp *InvoiceItem, err error) {
	r, err := stripe.request("GET", "invoiceitems/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
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
	if raw.Error != nil {
		// TODO: throw an error
	}
	return raw.Success, err
}

func (stripe *Stripe) ListInvoiceItems() (resp []*InvoiceItem, err error) {
	return stripe.QueryInvoiceItems(-1, -1, "")
}

func (stripe *Stripe) ListInvoiceItemsByCustomer(customer string) (resp []*InvoiceItem, err error) {
	return stripe.QueryInvoiceItems(-1, -1, customer)
}

func (stripe *Stripe) QueryInvoiceItems(count, offset int, customer string) (resp []*InvoiceItem, err error) {
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
		Count int            "count"
		Data  []*InvoiceItem "data"
		Error *RawError      "error"
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
