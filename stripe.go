package stripe

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	HOST    = "api.stripe.com"
	VERSION = 1
)

type Stripe struct {
	Host    string
	Version int
	AuthKey string
}

func New(auth string) *Stripe {
	return &Stripe{HOST, VERSION, auth}
}

type BadRequestError struct {
}

type UnauthorizedError struct {
}

type NotFoundError struct {
}

type ServerError struct {
}

type RawError struct {
	Type    string "type"
	Param   string "param"
	Code    string "code"
	Message string "message"
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

func (stripe *Stripe) request(method, url string, body string) (resp []byte, err error) {
	baseURL := "https://" + stripe.Host + "/v" + strconv.Itoa(stripe.Version) + "/"
	rbody := bytes.NewBufferString(body)
	req, err := http.NewRequest(method, baseURL+url, rbody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(stripe.AuthKey, "")
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
	resp, err = ioutil.ReadAll(hresp.Body)
	hresp.Body.Close()
	return resp, err
}
