package stripe

import (
        "testing"
        "io/ioutil"
)

func TestListInvoices(t *testing.T) {
        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        _, err = API.ListInvoices()
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
}
