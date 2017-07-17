package printorder

import (
	"testing"
	"net/url"
	"errors"
	"reflect"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"log"
	"fmt"
)

func TestCanShowPrintOrder(t *testing.T) {
	t.Parallel()
	id := 4
	r := Resource{Endpoint:SHOW,Id:id}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printorder.PrintOrder" {
			t.Errorf(`Expected a PrintOrder as model but got: "%s"`, iType)
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T: t,
		GetCall: cb,
	}

	_, err := Show(c, id)

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

func TestCanIndexPrintOrders(t *testing.T) {
	t.Parallel()
	r := Resource{Endpoint:INDEX}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printorder.IndexResponse" {
			t.Errorf(`Expected a PrintOrder as model but got: "%s"`, iType)
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

// Test helper Client Mock
type MockProductionAPIClient struct {
	ReturnError bool
	T *testing.T
	GetCall func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values))
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

func Example() {
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

	// Filter request to show only created delivery numbers after 2017-01-01 00:00:00.
	filter := common.QueryAttr(
		common.AttrQuery{
			Name: CREATED_AT,
			Value: "2017-07-15 00:00:00",
			Args: common.AttrArgs{
				Operator: common.OPERATOR_GREATER_EQUAL,
				Combinator: common.COMBINATOR_AND,
			},
		},
	)

	// Create with part of query to load relations with statuses and PrintData with files.
	with := common.QueryWith(WITH_STATUSES, WITH_PRINT_DATA_FILE)

	// Create scope query part.
	// Below filters PrintOrders to only return active and orders that has the status "exported".
	sc := []common.Scope{
		{
			Scope: "active",
		},
		{
			Scope: "exported",
		},
	}
	scope := common.QueryScope(sc)

	// Create a limit to fetch a maximum of 10 items
	limit := common.QueryLimit(10, 0)

	// Index PrintOrders
	po, err := Index(c, filter, scope, with, limit)
	if err != nil {
		log.Fatal(err)
	}

	// Prints out number of returned items in response.
	fmt.Printf("Total matches: %d, number of items in list: %d\n", po.Count, len(po.Data))
}

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
	po, err := Show(c, id)

	if err != nil {
		fmt.Printf("Could not find PrintOrder with id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find PrintOrder for id. The below would output true.
	fmt.Println(po.ID==id)
}

func ExampleIndex() {
	// See the package example for a more verbose example.

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

	// Index PrintOrders with limit
	po, err := Index(c, common.QueryLimit(10, 0))
	if err != nil {
		log.Fatal(err)
	}

	// Prints out number of returned items in response.
	fmt.Printf("Total matches: %d, number of items in list: %d\n", po.Count, len(po.Data))
}