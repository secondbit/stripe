package stripe

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	HOST     = "api.stripe.com"
	VERSION  = "v1"
	AUTH_KEY = "MY AUTH KEY"
)

type BadRequestError struct {
}

type UnauthorizedError struct {
}

type NotFoundError struct {
}

type ServerError struct {
}

type InvalidRequestError struct {
	Code    string
	Param   string
	Message string
}

type APIError struct {
	Code    string
	Param   string
	Message string
}

type CardError struct {
	Code    string
	Param   string
	Message string
}

func (c *CardError) Details() string {
	switch c.Code {
	case "invalid_number":
		return "The card number is invalid."
	case "incorrect_number":
		return "The card number is incorrect."
	case "invalid_expiry_month":
		return "The card's expiration month is invalid."
	case "invalid_expiry_year":
		return "The card's expiration year is invalid."
	case "invalid_cvc":
		return "The card's security code is invalid."
	case "expired_card":
		return "The card has expired."
	case "invalid_amount":
		return "An invalid amount was entered."
	case "incorrect_cvc":
		return "The card's security code is incorrect."
	case "card_declined":
		return "The card was declined."
	case "missing":
		return "There is no card on that customer."
	case "duplicate_transaction":
		return "A transaction with identical amount and credit card information was submitted very recently."
	case "processing_error":
		return "An error occurred while processing the card."
	}
	return ""
}

func apiRequest(method, url string, body string) (resp interface{}, err error) {
	baseURL := "https://" + HOST + "/" + VERSION + "/"
	rbody := bytes.NewBufferString(body)
	req, err := http.NewRequest(method, baseURL+url, rbody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(AUTH_KEY, "")
	hresp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch hresp.StatusCode {
	case 400:
		// TODO: Throw a BadRequest error
	case 401:
		// TODO: Throw an Unauthorized error
	case 402:
		// TODO: Throw a RequestFailed error
	case 404:
		// TODO: Throw a NotFound error
	case 500:
		// TODO: Throw a Server error
	case 502:
		// TODO: Throw a Server error
	case 503:
		// TODO: Throw a Server error
	case 504:
		// TODO: Throw a Server error
	case 200:
		// TODO: Keep calm, carry on
	default:
		// TODO: Throw a generic error
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(hresp.Body)
	hresp.Body.Close()
	err = json.Unmarshal(buf.Bytes(), &resp)
	return resp, err
}

type Charge struct {
	Amount      int
	Currency    string
	Card        *PartialCard
	Customer    string
	Description string
	Created     time.Time
	Fee         int
	ID          string
	LiveMode    bool
	Object      string
	Paid        bool
	Refunded    bool
}

func jsonToCharge(json interface{}) (resp *Charge, err error) {
	r := json.(map[string]interface{})
	created, err := time.Parse("%s", strconv.FormatInt(r["created"].(int64), 10))
	if err != nil {
		// TODO: throw an error
	}
	cardVal := r["card"].(map[string]interface{})
	card := PartialCard{
		Country:         cardVal["country"].(string),
		CVCCheck:        cardVal["cvc_check"].(string),
		AddressCheck:    cardVal["address_line1_check"].(string),
		AddressZipCheck: cardVal["address_zip_check"].(string),
		ExpMonth:        cardVal["exp_month"].(int),
		ExpYear:         cardVal["exp_year"].(int),
		LastFour:        cardVal["last4"].(string),
		Object:          cardVal["object"].(string),
		Type:            cardVal["type"].(string),
	}
	resp = &Charge{
		Amount:      r["amount"].(int),
		Currency:    r["currency"].(string),
		Customer:    r["customer"].(string),
		Description: r["description"].(string),
		Created:     created,
		Fee:         r["fee"].(int),
		ID:          r["id"].(string),
		LiveMode:    r["livemode"].(bool),
		Object:      r["object"].(string),
		Paid:        r["paid"].(bool),
		Refunded:    r["refunded"].(bool),
		Card:        &card,
	}
	return
}

func CreateCharge(amount int, currency string, customer *Customer, card *Card, description string) (resp *Charge, err error) {
	if currency != "usd" {
		// TODO: throw an error
	}
	if customer != nil && card != nil {
		// TODO: throw an error
	}
	values := make(url.Values)
	values.Set("amount", strconv.Itoa(amount))
	values.Set("currency", currency)
	if description != "" {
		values.Set("description", description)
	}
	if customer != nil {
		values.Set("customer", customer.ID)
	}
	if card != nil {
		if card.Token != nil {
			values.Set("card", card.Token.ID)
		} else {
			if card.Number == "" {
				// TODO: throw an error
			} else if card.ExpMonth < 0 {
				// TODO: throw an error
			} else if card.ExpYear < 0 {
				// TODO: throw an error
			} else {
				values.Set("card[\"number\"]", card.Number)
				values.Set("card[\"exp_month\"]", strconv.Itoa(card.ExpMonth))
				values.Set("card[\"exp_year\"]", strconv.Itoa(card.ExpYear))
				if card.CVC != "" {
					values.Set("card[\"cvc\"]", card.CVC)
				}
				if card.Name != "" {
					values.Set("card[\"name\"]", card.Name)
				}
				if card.Address1 != "" {
					values.Set("card[\"address_line1\"]", card.Address1)
				}
				if card.Address2 != "" {
					values.Set("card[\"address_line2\"]", card.Address2)
				}
				if card.AddressZip != "" {
					values.Set("card[\"address_zip\"]", card.AddressZip)
				}
				if card.AddressState != "" {
					values.Set("card[\"address_state\"]", card.AddressState)
				}
				if card.AddressCountry != "" {
					values.Set("card[\"address_country\"]", card.AddressCountry)
				}
			}
		}
	}
	body := values.Encode()
	r, err := apiRequest("POST", "charges", body)
	if err != nil {
		return nil, err
	}
	_, fail := r.(map[string]interface{})["error"]
	if fail {
		// TODO: Throw an error
	}
	resp, err = jsonToCharge(r)
	return
}

func GetCharge(id string) (resp *Charge, err error) {
	if id == "" {
		// TODO: throw an error
	}
	r, err := apiRequest("GET", "charges/"+id, "")
	if err != nil {
		return nil, err
	}
	_, fail := r.(map[string]interface{})["error"]
	if fail {
		// TODO: Throw an error
	}
	resp, err = jsonToCharge(r)
	return
}

func (c *Charge) Refund(amount int) (resp *Charge, err error) {
	if c == nil {
		// TODO: throw an error
	}
	var body string
	if amount >= 0 {
		values := make(url.Values)
		values.Set("amount", strconv.Itoa(amount))
		body = values.Encode()
	}
	r, err := apiRequest("POST", "charges/"+c.ID+"/refund", body)
	if err != nil {
		return nil, err
	}
	_, fail := r.(map[string]interface{})["error"]
	if fail {
		// TODO: Throw an error
	}
	resp, err = jsonToCharge(r)
	return
}

func ListCharges(count, offset int, customer string) (resp []*Charge, err error) {
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
	r, err := apiRequest("GET", "charges"+params, "")
	if err != nil {
		return nil, err
	}
	_, fail := r.(map[string]interface{})["error"]
	if fail {
		// TODO: Throw an error
	}
	resp = []*Charge{}
	for _, charge := range r.(map[string][]interface{})["data"] {
		c, err := jsonToCharge(charge)
		if err != nil {
			// TODO: throw an error
		}
		resp = append(resp, c)
	}
	return
}

type Customer struct {
	Object      string
	Description string
	LiveMode    bool
	Created     time.Time
	ActiveCard  *PartialCard
	ID          string
}

func CreateCustomer(card *Card, coupon *Coupon, email, description string, plan *Plan, trial_end time.Time) (resp *Customer, err error)

func GetCustomer(id string) (resp *Customer, err error)

func (c *Customer) Update(card *Card, coupon *Coupon, email, description string) (resp *Customer, err error)

func (c *Customer) Delete() (success bool, err error)

func ListCustomers(count, offset int) (resp []*Customer, err error)

type Card struct {
	Token          *Token
	Number         string
	ExpMonth       int
	ExpYear        int
	CVC            string
	Name           string
	Address1       string
	Address2       string
	AddressZip     string
	AddressState   string
	AddressCountry string
}

type PartialCard struct {
	Country         string
	CVCCheck        string
	AddressCheck    string
	AddressZipCheck string
	ExpMonth        int
	ExpYear         int
	LastFour        string
	Object          string
	Type            string
}

type Token struct {
	Created  time.Time
	Currency string
	Used     bool
	Amount   int
	Object   string
	LiveMode bool
	ID       string
	Card     *PartialCard
}

func CreateCardToken(card *Card, amount int, currency string) (resp *Token, err error)

func GetCardToken(id string) (resp *Token, err error)

type Plan struct {
	Object   string
	Amount   int
	Interval string
	Name     string
	Currency string
	ID       string
}

type Subscription struct {
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	Status             string
	Plan               *Plan
	Object             string
	TrialStart         time.Time
	TrialEnd           time.Time
	Customer           string
}

func CreatePlan(id string, amount int, currency, interval, name string, trial_period_days int) (resp *Plan, err error)

func GetPlan(id string) (resp *Plan, err error)

func (p *Plan) Update(name string) (resp *Plan, err error)

func (p *Plan) Delete() (success bool, err error)

func ListPlans(count, offset int) (resp []*Plan, err error)

func (c *Customer) Subscribe(plan, coupon string, prorate bool, trial_end time.Time, card *Card) (resp *Subscription, err error)

func (c *Customer) Unsubscribe(at_period_end bool) (resp *Subscription, err error)

type InvoiceItem struct {
	Date        time.Time
	Description string
	Currency    string
	Amount      int
	Object      string
	ID          string
}

func CreateInvoiceItem(customer *Customer, amount int, currency, invoice, description string) (resp *InvoiceItem, err error)

func GetInvoiceItem(id string) (resp *InvoiceItem, err error)

func (i *InvoiceItem) Update(amount int, description string) (resp *InvoiceItem, err error)

func (i *InvoiceItem) Delete() (success bool, err error)

func ListInvoiceItems(count, offset int, customer string) (resp []*InvoiceItem, err error)

type Invoice struct {
	Subtotal int
	Total    int
	Lines    map[string][]interface{}
	Object   string
}

func GetInvoice(id string) (resp *Invoice, err error)

func ListInvoices(count, offset int, customer string) (resp []*Invoice, err error)

func (c *Customer) NextInvoice() (resp *Invoice, err error)

type Coupon struct {
	PercentOff       int
	Duration         string
	DurationInMonths int
	Object           string
	ID               string
}

func CreateCoupon(id string, percent_off int, duration string, duration_in_months, max_redemption int, redeem_by time.Time) (resp *Coupon, err error)

func GetCoupon(id string) (resp *Coupon, err error)

func (c *Coupon) Delete() (success bool, err error)

func ListAllCoupons(count, offset int) (resp []*Coupon, err error)

type Event struct {
	Type            string
	Created         time.Time
	PendingWebhooks string
	Data            map[string]interface{}
	LiveMode        bool
	ID              string
}

func GetEvent(id string) (resp *Event, err error)

func ListEvents(count, offset int, eventType string, created interface{})
