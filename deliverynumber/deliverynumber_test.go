package deliverynumber_test

import (
	"testing"
	. "github.com/publitsweden/ProductionAPIGoSDK/deliverynumber"
	"reflect"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/url"
	"errors"
	"net/http"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"log"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
)

func TestResourceImplementsEndpointer(t *testing.T) {
	t.Parallel()

	r := Resource{}
	ei := reflect.TypeOf((*production.Endpointer)(nil)).Elem()

	if !reflect.TypeOf(r).Implements(ei) {
		t.Error("Resource does not implement production.Endpointer but was epected to.")
	}
}

func TestCanIndexDeliveryNumbers(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint:INDEX}
	c := &MockProductionAPIClient{}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "deliverynumber.IndexResponse" {
			t.Errorf(`Expected a Delivery number index response as model but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c.GetCall = cb

	qp := func(q url.Values){}
	_, err := Index(c, qp)

	if err != nil {
		t.Error("Received an error but was not expecting one.")
	}
}

func TestCanShowDeliveryNumber(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint:SHOW, Id: 9}
	c := &MockProductionAPIClient{}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "deliverynumber.DeliveryNumber" {
			t.Errorf(`Expected a Delivery number as model but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c.GetCall = cb

	qp := func(q url.Values){}
	_, err := Show(c, r.Id, qp)

	if err != nil {
		t.Error("Received an error but was not expecting one.")
	}
}

func TestCanStoreDeliveryNumber(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint: POST}

	s := New(1, "ABS3423", "somemessage")

	c := &MockProductionAPIClient{}

	c.PostPutDeleteCall = func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		if payload != s {
			t.Error("Provided payload did not match expected.")
		}

		// Since result is interface{} we must assert it's kind and value using reflect.
		switch reflect.TypeOf(reflect.ValueOf(result).Elem().Interface()).Kind() {
		case reflect.Slice:
			sl := reflect.ValueOf(reflect.ValueOf(result).Elem().Interface())
			v := sl.Index(0).Interface()
			if v != s {
				t.Error("Provided result slice did not contain pointer to payload.")
			}

		default:
			t.Error("Provided result was not a slice, but was expected to be.")
		}
	}

	err := s.Store(c)

	if err != nil {
		t.Error("Got error did not expect one.", err.Error())
	}
}

func TestCanNotStoreExistingDeliveryNumber(t *testing.T) {
	t.Parallel()

	s := New(1, "ABS3423", "somemessage")
	s.ID = 1

	c := &MockProductionAPIClient{}
	err := s.Store(c)

	if err == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

func TestCanUpdateDeliveryNumber(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint: PUT, Id: 5}

	s := New(1, "ABS3423", "somemessage")
	s.ID = r.Id

	c := &MockProductionAPIClient{}

	c.PostPutDeleteCall = func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		if payload != s {
			t.Error("Provided payload did not match expected.")
		}

		// We want updates to be made to the same model we sent in. Thus same result.
		if result != s {
			t.Error("Provided result was not a slice, but was expected to be.")
		}
	}

	err := s.Update(c)

	if err != nil {
		t.Error("Got error did not expect one.", err.Error())
	}
}

func TestCanNotUpdateNewDeliveryNumber(t *testing.T) {
	t.Parallel()

	s := New(1, "ABS3423", "somemessage")

	c := &MockProductionAPIClient{}

	err := s.Update(c)

	if err == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

func TestCanDeleteDeliveryNumber(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint: DELETE, Id: 5}

	s := &DeliveryNumber{}
	s.ID = r.Id

	c := &MockProductionAPIClient{}
	c.PostPutDeleteCall = func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		// We want updates to be made to the same model we sent in. Thus same result.
		if result != s {
			t.Error("Provided result was not a slice, but was expected to be.")
		}
	}

	err := s.Delete(c)
	if err != nil {
		t.Error("Received an error but was not expecting one.", err.Error())
	}
}

func TestCanNotDeleteNonExistingDeliveryNumber(t *testing.T) {
	t.Parallel()

	s := &DeliveryNumber{}

	c := &MockProductionAPIClient{}
	err := s.Delete(c)
	if err == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}


// Test helper Client Mock
type MockProductionAPIClient struct {
	ReturnError bool
	T *testing.T
	GetCall func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values))
	PostPutDeleteCall func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(q *http.Header))
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

func (c *MockProductionAPIClient) Post(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	if c.ReturnError {
		return errors.New("Some error")
	}

	if c.PostPutDeleteCall != nil {
		c.PostPutDeleteCall(c.T, endpoint, payload, result, headers...)
	}

	return nil
}

func (c *MockProductionAPIClient) Put(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	if c.ReturnError {
		return errors.New("Some error")
	}

	if c.PostPutDeleteCall != nil {
		c.PostPutDeleteCall(c.T, endpoint, payload, result, headers...)
	}

	return nil
}

func (c *MockProductionAPIClient) Delete(endpoint production.Endpointer, result interface{}, headers ...func(h *http.Header)) error {
	if c.ReturnError {
		return errors.New("Some error")
	}

	if c.PostPutDeleteCall != nil {
		c.PostPutDeleteCall(c.T, endpoint, nil, result, headers...)
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
			Value: "2017-01-01 00:00:00",
			Args: common.AttrArgs{
				Operator: common.OPERATOR_GREATER_EQUAL,
				Combinator: common.COMBINATOR_AND,
			},
		},
	)

	// Create a limit to fetch a maximum of 10 items
	limit := common.QueryLimit(10, 0)

	// Run Index to retrieve list of DeliveryNumbers.
	d, err := Index(c, filter, limit)

	if err != nil {
		fmt.Printf("Could not list delivery numbers after 2017-01-01. Got errors: %v\n", err.Error())
	}

	// Print the number of DeliveryNumbers in list from response.
	fmt.Println(len(d.Data))
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
	d, err := Show(c, id)

	if err != nil {
		fmt.Printf("Could not find id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find delivery number. The below would output true.
	fmt.Println(d.ID==id)
}

func ExampleDeliveryNumber_Store() {
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

	d := New(1, "somedeliverynumber", "Some message")

	err := d.Store(c)
	if err != nil {
		fmt.Println("Could not store delivery number.", err.Error())
	}

	// If storing was succesful the created at date (and other attributes) will be updated based on response from Publit.
	fmt.Println(d.CreatedAt)
}

func ExampleDeliveryNumber_Update() {
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

	// Create an "existing" delivery number to be able to perform an update.
	// An existing delivery number has an ID.
	d := &DeliveryNumber{
		ID: 1,
		Message: "Some new message for delivery number with id 1",
	}

	err := d.Update(c)
	if err != nil {
		log.Println("Could not store delivery number.", err.Error())
	}

	updated,_ := d.UpdatedAt.ConvertPublitTimeToTime()
	created,_ := d.CreatedAt.ConvertPublitTimeToTime()
	// If update was succesful the updatedAt timestamp should have been updated.
	fmt.Println(updated.After(created))
}

func ExampleDeliveryNumber_Delete() {
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

	// Create an "existing" delivery number to be able to perform DELETE.
	// An existing delivery number has an ID.
	id := 1
	d := &DeliveryNumber{ID: id}

	err := d.Delete(c)
	if err != nil {
		fmt.Println("Could not delete delivery number.", err.Error())
	}

	// If trying to show the removed delivery number the response will be 404 NotFount.
	_, err = Show(c, id)

	if err != nil {
		fmt.Println("Could not find delivery number. Must be removed...")
	}
}