package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Event struct {
	Type string "type"
	Data struct {
		Object interface{} "object"
	}
	PendingWebhooks int    "pending_webhooks"
	LiveMode        bool   "livemode"
	Created         int    "created"
	ID              string "id"
	Object          string "object"
}

func (stripe *Stripe) GetEvent(id string) (resp *Event, err error) {
	r, err := stripe.request("GET", "events/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	return
}

func (stripe *Stripe) ListEvents() (resp []*Event, err error) {
	return stripe.QueryEvents("", -1, -1, "", "")
}

func (stripe *Stripe) ListEventsByType(event_type string) (resp []*Event, err error) {
	return stripe.QueryEvents(event_type, -1, -1, "", "")
}

func (stripe *Stripe) ListEventsOnDate(date string) (resp []*Event, err error) {
	return stripe.QueryEvents("", -1, -1, date, "")
}

func (stripe *Stripe) ListEventsComparedToDate(comparison, date string) (resp []*Event, err error) {
	return stripe.QueryEvents("", -1, -1, date, comparison)
}

func (stripe *Stripe) QueryEvents(event_type string, count, offset int, date, comparison string) (resp []*Event, err error) {
	values := make(url.Values)
	if count >= 0 {
		values.Set("count", strconv.Itoa(count))
	}
	if offset >= 0 {
		values.Set("offset", strconv.Itoa(offset))
	}
	if event_type != "" {
		values.Set("type", event_type)
	}
	if date != "" && comparison != "" {
		values.Set("created["+comparison+"]", date)
	} else if date != "" {
		values.Set("created", date)
	}
	params := values.Encode()
	if params != "" {
		params = "?" + params
	}
	r, err := stripe.request("GET", "events"+params, "")
	if err != nil {
		return nil, err
	}
	var raw struct {
		Count int       "count"
		Data  []*Event  "data"
		Error *RawError "error"
	}
	err = json.Unmarshal(r, &raw)
	if err != nil {
		return nil, err
	}
	if raw.Error.Code != "" {
		// TODO: throw an error
	}
	resp = raw.Data
	return
}
