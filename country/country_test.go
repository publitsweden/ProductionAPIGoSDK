package country_test

import (
	"net/url"
	"testing"
	"github.com/publitsweden/ProductionAPIGoSDK"
	. "github.com/publitsweden/ProductionAPIGoSDK/country"
	"errors"
	"reflect"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"log"
)

func TestCanIndexCountries(t *testing.T) {
	t.Parallel()
	r := Resource{Endpoint:INDEX}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "country.IndexResponse" {
			t.Errorf(`Expected an IndexResponse as model but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T: t,
		GetCall: cb,
	}

	qp := func(q url.Values) {}
	_, err := Index(c, qp)

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

func TestCanShowCountries(t *testing.T) {
	t.Parallel()
	r := Resource{Endpoint:SHOW,Id:4}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "country.Country" {
			t.Errorf(`Expected an IndexResponse as model but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T: t,
		GetCall: cb,
	}

	qp := func(q url.Values) {}
	_, err := Show(c, r.Id, qp)

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

// Test helper Client Mock
type MockProductionAPIClient struct {
	ReturnError bool
	T           *testing.T
	GetCall    func(t *testing.T, endpoint production.Endpointer, payload interface{}, queryParams ...func(q url.Values))
}

func (c *MockProductionAPIClient) Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) (error) {
	if c.ReturnError {
		return errors.New("Some error")
	}

	if c.GetCall != nil {
		c.GetCall(c.T, endpoint, model, queryParams...)
	}

	return nil
}

// Examples
func ExampleShow() {
	// Create APIClient to perform calls
	c := production.APIClient{
		Client: client.New(
			func(c *client.Client) {
				c.User = "MyUser"
				c.Password = "MyPassword"
			},
		),
		BaseUrl: "https://url.to.api",
	}

	id := 1
	country, err := Show(c, id)

	if err != nil {
		fmt.Printf("Could not find Status with id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find country for id. The below would output true.
	fmt.Println(country.ID==id)
}

func ExampleIndex() {
	// Create APIClient to perform calls
	c := production.APIClient{
		Client: client.New(
			func(c *client.Client) {
				c.User = "MyUser"
				c.Password = "MyPassword"
			},
		),
		BaseUrl: "https://url.to.api",
	}

	// Index Countries with limit
	resp, err := Index(c, common.QueryLimit(10, 0))
	if err != nil {
		log.Fatal(err)
	}

	// Prints out number of returned items in response.
	fmt.Printf("Total matches: %d, number of items in list: %d\n", resp.Count, len(resp.Data))
}
