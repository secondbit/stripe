package stripe

import (
	"bytes"
        "fmt"
	"io/ioutil"
	"net/http"
        "net/url"
	"strconv"
        "strings"
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
	return &Stripe{HOST, VERSION, strings.TrimSpace(auth)}
}

type BadRequestError struct {
        Message string "message"
        Request *http.Request "request"
}

func (err *BadRequestError) Error() string {
        if err.Request != nil {
                err.Message = fmt.Sprintf("%v\nRequest: %v", err.Message, err.Request)
        }
        return err.Message
}

func ThrowBadRequest(req *http.Request) *BadRequestError {
        return &BadRequestError {
                Message: "Error: Bad request.",
                Request: req,
        }
}

type UnauthorizedError struct {
        Message string "message"
        Auth string "auth"
}

func (err *UnauthorizedError) Error() string {
        if err.Auth == "" {
                return "No authorization set."
        }
        return fmt.Sprintf("%v\nAuth: %v", err.Message, err.Auth)
}

func ThrowUnauthorized(req *http.Request) *UnauthorizedError {
        return &UnauthorizedError {
                Message: "Error: Unauthorized.",
                Auth: req.Header["Authorization"][0],
        }
}

type RequestFailedError struct {
        Message string "message"
        Request *http.Request "request"
}

func (err *RequestFailedError) Error() string {
        if err.Request != nil {
                err.Message = fmt.Sprintf("%v\nRequest: %v", err.Message, err.Request)
        }
        return err.Message
}

func ThrowRequestFailed(req *http.Request) *RequestFailedError {
        return &RequestFailedError {
                Message: "Error: Request failed.",
                Request: req,
        }
}

type NotFoundError struct {
        Message string "message"
        URL *url.URL "url"
}

func (err *NotFoundError) Error() string {
        if err.URL != nil {
                err.Message = fmt.Sprintf("%v\nURL: %v", err.Message, err.URL)
        }
        return err.Message
}

func ThrowNotFound(req *http.Request) *NotFoundError {
        return &NotFoundError {
                Message: "Error: Resource not found.",
                URL: req.URL,
        }
}

type ServerError struct {
        Message string "message"
}

func (err *ServerError) Error() string {
        return err.Message
}

func ThrowServer(req *http.Request) *ServerError {
        return &ServerError {
                Message: "Error: Internal Server Error.",
        }
}

type UnknownError struct {
        Message string "message"
}

func (err *UnknownError) Error() string {
        return err.Message
}

func ThrowUnknown(code int) *UnknownError {
        return &UnknownError {
                Message: "Error: Unknown error occurred.\nStatus Code: "+strconv.Itoa(code),
        }
}

type RawError struct {
	Type    string "type"
	Param   string "param"
	Code    string "code"
	Message string "message"
}

func (err *RawError) Error() string {
        return fmt.Sprintf("Error (%v): %v\nType: %v\nParam: %v", err.Code, err.Message, err.Type, err.Param)
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
	req.SetBasicAuth(stripe.AuthKey,"")
	hresp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch hresp.StatusCode {
	case 400:
		return nil, ThrowBadRequest(req)
	case 401:
                return nil, ThrowUnauthorized(req)
	case 402:
		return nil, ThrowRequestFailed(req)
	case 404:
                return nil, ThrowNotFound(req)
	case 500:
                return nil, ThrowServer(req)
	case 502:
                return nil, ThrowServer(req)
	case 503:
                return nil, ThrowServer(req)
	case 504:
                return nil, ThrowServer(req)
	case 200:
        	resp, err = ioutil.ReadAll(hresp.Body)
	        hresp.Body.Close()
        	return resp, err
	default:
		return nil, ThrowUnknown(hresp.StatusCode)
	}
        return nil, ThrowUnknown(hresp.StatusCode)
}
