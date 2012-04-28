package stripe

import (
	"io/ioutil"
	"testing"
)

//TODO: TestGetCoupons
//TODO: TestCreateCoupon
//TODO: TestDeleteCoupon

func TestListCoupons(t *testing.T) {
	key, err := ioutil.ReadFile("key")
	if err != nil {
		t.Fatalf("err = %v, want %v", err, nil)
	}
	API := New(string(key))
	_, err = API.ListCoupons()
	if err != nil {
		t.Fatalf("err = %v, want %v", err, nil)
	}
}
