package stripe

import (
	"io/ioutil"
	"testing"
)

// TODO: TestCreatePlan
// TODO: TestReadPlan
// TODO: TestUpdatePlan
// TODO: TestDeletePlan

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
