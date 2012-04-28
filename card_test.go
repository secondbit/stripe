package stripe

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	VALID = &Card{
		Number:         "4242424242424242",
		ExpMonth:       3,
		ExpYear:        time.Now().Year() + 4,
		CVC:            "123",
		Name:           "Oso de Peluche",
		AddressLine1:   "123 Awesome Street",
		AddressLine2:   "Apartment 5",
		Zip:            "12345",
		State:          "",
		AddressCountry: "Spain",
	}
	INVALID = &Card{
		Number:         "2",
		ExpMonth:       3,
		ExpYear:        time.Now().Year() - 1,
		CVC:            "-1",
		Name:           "Oso de Peluche",     // TODO: Come up with an invalid value
		AddressLine1:   "123 Awesome Street", // TODO: Come up with an invalid value
		AddressLine2:   "Apartment 5",        // TODO: Come up with an invalid value
		Zip:            "12345",              // TODO: Come up with an invalid value
		State:          "",                   // TODO: Come up with an invalid value
		AddressCountry: "Spain",              // TODO: come up with an invalid value
	}
	key, err = ioutil.ReadFile("key")
)

func TestCreateCardToken(t *testing.T) {
	if err != nil {
		t.Fatalf("err = %v, want %v", err, nil)
	}
	API := New(string(key))
	var token *Token
	token, err = API.CreateToken(VALID)
	if err != nil {
		t.Fatalf("err = %v, want %v", err, nil)
	}
	if token == nil {
		t.Fatalf("token is nil, should be set")
	}
	t.Logf("token = %v", token)
	if token.Card.ExpYear != VALID.ExpYear {
		t.Errorf("ExpYear is %v, expected %v", token.Card.ExpYear, VALID.ExpYear)
	}
	if token.Card.ExpMonth != VALID.ExpMonth {
		t.Errorf("ExpMonth is %v, expected %v", strconv.Itoa(token.Card.ExpMonth), VALID.ExpMonth)
	}
	if !strings.HasSuffix(VALID.Number, token.Card.LastFour) {
		t.Errorf("token.Card.LastFour is %v, expected %v%v%v%v", token.Card.LastFour, VALID.Number[len(VALID.Number)-1], VALID.Number[len(VALID.Number)-2], VALID.Number[len(VALID.Number)-3], VALID.Number[len(VALID.Number)-4])
	}
	if token.Card.Name != VALID.Name {
		t.Errorf("token.Card.Name is %v, expected %v", token.Card.Name, VALID.Name)
	}
	if token.Card.AddressCountry != VALID.AddressCountry {
		t.Error("token.Card.AddressCountry is %v, expected %v", token.Card.AddressCountry, VALID.AddressCountry)
	}
	if token.Card.AddressLine1 != VALID.AddressLine1 {
		t.Error("token.Card.AddressLine1 is %v, expected %v", token.Card.AddressLine1, VALID.AddressLine1)
	}
	if token.Card.AddressLine2 != VALID.AddressLine2 {
		t.Error("token.Card.AddressLine2 is %v, expected %v", token.Card.AddressLine2, VALID.AddressLine2)
	}
	if token.Card.Zip != VALID.Zip {
		t.Error("token.Card.Zip is %v, expected %v", token.Card.Zip, VALID.Zip)
	}
	if token.Card.State != VALID.State {
		t.Error("token.Card.State is %v, expected %v", token.Card.State, VALID.State)
	}
}

// TODO: Test with every permutation of values

// TODO: Update this test.

/*func TestGetCardToken(t *testing.T) {
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        API := New(string(key))
        var token, token2 *CardToken
        token, err = API.CreateCardTokenWithAll(VALID.Number, VALID.ExpMonth, VALID.ExpYear, VALID.CVC, VALID.Name, VALID.Address1, VALID.Address2, VALID.Zip, VALID.State, VALID.Country)
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        if token == nil {
                t.Fatalf("token is nil, should be set")
        }
        if token.ID == "" {
                t.Fatalf("token.ID not set")
        }
        token2, err = API.GetCardToken(token.ID)
        if err != nil {
                t.Fatalf("err = %v, want %v", err, nil)
        }
        if token2 == nil {
                t.Fatalf("token2 is nil, should be set")
        }
        if token.Card.ExpYear != token2.Card.ExpYear {
                t.Errorf("ExpYear is %v, expected %v", token2.Card.ExpYear, token.Card.ExpYear)
        }
        if token.Card.ExpMonth != token2.Card.ExpMonth {
                t.Errorf("ExpMonth is %v, expected %v", token2.Card.ExpMonth, token.Card.ExpMonth)
        }
        if token.Card.LastFour != token2.Card.LastFour {
                t.Errorf("token2.Card.LastFour is %v, expected %v", token2.Card.LastFour, token.Card.LastFour)
        }
        if token.Card.Name != token2.Card.Name {
                t.Errorf("token2.Card.Name is %v, expected %v", token2.Card.Name, token.Card.Name)
        }
        if token.Card.AddressCountry != token2.Card.AddressCountry {
                t.Error("token2.Card.AddressCountry is %v, expected %v", token2.Card.AddressCountry, token.Card.Country)
        }
        if token.Card.Address1 != token2.Card.Address1 {
                t.Error("token2.Card.Address1 is %v, expected %v", token2.Card.Address1, token.Card.Address1)
        }
        if token.Card.Address2 != token2.Card.Address2 {
                t.Error("token2.Card.Address2 is %v, expected %v", token2.Card.Address2, token.Card.Address2)
        }
        if token.Card.Zip != token2.Card.Zip {
                t.Error("token2.Card.Zip is %v, expected %v", token2.Card.Zip, token.Card.Zip)
        }
        if token.Card.State != token2.Card.State {
                t.Error("token2.Card.State is %v, expected %v", token2.Card.State, token.Card.State)
        }

}*/
