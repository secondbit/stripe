package stripe

import (
        "testing"
        "io/ioutil"
)

func TestListCharges(t *testing.T) {
        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        t.Logf("key = %v", string(key))
        API := New(string(key))
        t.Logf("API.AuthKey = '%v'", API.AuthKey)
        _, err = API.ListCharges()
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
}
