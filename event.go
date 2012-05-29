package stripe

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Event struct {
        Type string `json:"type"`
	Data struct {
                Object interface{} `json:"object"`
	}
        PendingWebhooks int       `json:"pending_webhooks"`
        LiveMode        bool      `json:"livemode"`
        Created         int64       `json:"created"`
        ID              string    `json:"id"`
        Object          string    `json:"object"`
        Error           *RawError `json:"error"`
}

func (stripe *Stripe) GetEvent(id string) (resp *Event, err error) {
	r, err := stripe.request("GET", "events/"+id, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		//TODO: Throw an error
	}
	return
}

func (stripe *Stripe) ListEvents(event *Event, count, offset int, comparison string ) (resp []*Event, err error) {
	values := make(url.Values)
	if count >= 0 {
		values.Set("count", strconv.Itoa(count))
	}
	if offset >= 0 {
		values.Set("offset", strconv.Itoa(offset))
	}
        if event != nil {
        	if event.Type != "" {
	        	values.Set("type", event.Type)
        	}
	        if event.Created > 0 && comparison != "" {
                        if comparison != "e" {
        		        values.Set("created["+comparison+"]", strconv.FormatInt(event.Created, 10))
                        } else {
                                values.Set("created", strconv.FormatInt(event.Created, 10))
                        }
        	}
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
	if raw.Error != nil {
		// TODO: throw an error
	}
	resp = raw.Data
	return
}
