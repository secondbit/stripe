package stripe

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
}

func apiRequest(method, url string, body io.Reader) (resp interface{}, err error) {
	baseURL := "https://" + HOST + "/" + VERSION + "/"
	req := http.NewRequest(method, baseURL+url, body)
	req.setBasicAuth(AUTH_KEY, "")
	hresp, err = http.DefaultClient.Do(req)
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
	var resp interface{}
	err = json.Unmarshal(r.Body, &resp)
	hresp.Body.Close()
	return resp, err
}

type Charge struct {
	Amount   int
	Currency string
	Card     struct {
		Country  string
		CVCCheck string
		ExpMonth int
		ExpYear  int
		LastFour string
		Object   string
		Type     string
	}
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

func jsonToCharge(json map[string]interface{}) (resp *Charge, err error) {
	created, err := time.Parse("%s", strconv.FormatInt(r["created"], 10))
	if err != nil {
		// TODO: throw an error
	}
	card := PartialCard{
		Country:         r["card"]["country"],
		CVCCheck:        r["card"]["cvc_check"],
		AddressCheck:    r["card"]["address_line1_check"],
		AddressZipCheck: r["card"]["address_zip_check"],
		ExpMonth:        r["card"]["exp_month"],
		ExpYear:         r["card"]["exp_year"],
		LastFour:        r["card"]["last4"],
		Object:          r["card"]["object"],
		Type:            r["card"]["type"],
	}
	resp := Charge{
		Amount:      r["amount"],
		Currency:    r["currency"],
		Customer:    r["customer"],
		Description: r["description"],
		Created:     created,
		Fee:         r["fee"],
		ID:          r["id"],
		LiveMode:    r["livemode"],
		Object:      r["object"],
		Paid:        r["paid"],
		Refunded:    r["refunded"],
		Card:        card,
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
	values.Set("amount", amount)
	values.Set("currency", currency)
	if description != nil {
		values.Set("description", description)
	}
	if customer != nil {
		values.Set("customer", customer.ID)
	}
	if card != nil {
		if card.Token != nil {
			values.set("card", card.Token.ID)
		} else {
			if card.Number == nil {
				// TODO: throw an error
			} else if card.ExpMonth == nil {
				// TODO: throw an error
			} else if card.ExpYear == nil {
				// TODO: throw an error
			} else {
				values.Set("card[\"number\"]", card.Number)
				values.Set("card[\"exp_month\"]", card.ExpMonth)
				values.Set("card[\"exp_year\"]", card.ExpYear)
				if card.CVC != nil {
					values.Set("card[\"cvc\"]", card.CVC)
				}
				if card.Name != nil {
					values.Set("card[\"name\"]", card.Name)
				}
				if card.Address1 != nil {
					values.Set("card[\"address_line1\"]", card.Address1)
				}
				if card.Address2 != nil {
					values.Set("card[\"address_line2\"]", card.Address2)
				}
				if card.AddressZip != nil {
					values.Set("card[\"address_zip\"]", card.AddressZip)
				}
				if card.AddressState != nil {
					values.Set("card[\"address_state\"]", card.AddressState)
				}
				if card.AddressCountry != nil {
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
	_, fail := r["error"]
	if fail {
		// TODO: Throw an error
	}
	resp, err = jsonToCharge(r)
	return
}

func GetCharge(id string) (resp *Charge, err error) {
	if id == nil {
		// TODO: throw an error
	}
	r, err := apiRequest("GET", "charges/"+id, nil)
	if err != nil {
		return nil, err
	}
	_, fail := r["error"]
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
	if amount != nil {
		values := make(url.Values)
		values.Set("amount", amount)
		body = values.Encode()
	}
	r, err := apiRequest("POST", "charges/"+id+"/refund", body)
	if err != nil {
		return nil, err
	}
	_, fail := r["error"]
	if fail {
		// TODO: Throw an error
	}
	resp, err = jsonToCharge(r)
	return
}

func ListCharges(count, offset int, customer string) (resp []*Charge, err error) {
	values := make(url.Values)
	if count != nil {
		values.Set("count", count)
	}
	if offset != nil {
		values.Set("offset", offset)
	}
	if customer != nil {
		values.Set("customer", customer)
	}
	params = values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := apiRequest("GET", "charges"+params, nil)
	if err != nil {
		return nil, err
	}
	_, fail := r["error"]
	if fail {
		// TODO: Throw an error
	}
	resp = []*Charge{}
	for _, charge := range r["data"] {
		c, err = jsonToCharge(r)
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
	ExpMonth       string
	ExpYear        string
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
