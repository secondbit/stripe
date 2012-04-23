package stripe

import (
        "testing"
        "io/ioutil"
        "time"
)

// TODO: TestGetCustomer
// TODO: TestUpdateCustomer
// TODO: TestDeleteCustomer

func TestListCustomers(t *testing.T) {
        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        _, err = API.ListCustomers()
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
}

func TestCreateCustomerWithToken(t *testing.T) {
        desc := "Test customer for Stripe's Go bindings."
        email := "test@gostripe.com"
        /*token := "mytesttoken"
        plan := "supertestplan"
        trial_end := time.Now()
        coupon := "TEST25OFF"*/

        key, err := ioutil.ReadFile("key")
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        //customer, err := API.CreateCustomerWithToken(email, token, desc, plan, trial_end, coupon)
        customer, err := API.CreateCustomerWithToken(email, "", desc, "", time.Now(), "")
        if err != nil {
                t.Errorf("err = %v, want %v", err, nil)
        }
        if customer == nil {
                t.Error("customer is nil, should be populated.")
        }
        if customer.Description != desc {
                t.Error("customer.Description = %v, want %v", customer.Description, desc)
        }
}
