package stripe

import (
        "testing"
        "io/ioutil"
)

func TestListPlans(t *testing.T) {
        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        _, err = API.ListPlans()
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
}
