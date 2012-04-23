package stripe

import (
        "testing"
        "io/ioutil"
)

//TODO: TestGetEvent

func TestListEvents(t *testing.T) {
        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        _, err = API.ListEvents()
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
}
