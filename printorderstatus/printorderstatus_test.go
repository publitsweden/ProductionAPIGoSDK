package printorderstatus

import (
	"errors"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/http"
	"reflect"
	"testing"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"net/url"
	"log"
)

func TestCanGetLastStatus(t *testing.T) {
	t.Parallel()
	lastTime := common.PublitTime("2017-03-01 00:00:02")

	list := StatusList{
		{
			UpdatedAt: "2017-01-01 00:00:00",
		},
		{
			UpdatedAt: "2015-01-02 00:00:00",
		},
		// This should be last
		{
			UpdatedAt: lastTime,
		},
		{
			UpdatedAt: "2017-03-01 00:00:00",
		},
	}

	l := list.GetLast()

	t1, _ := l.UpdatedAt.ConvertPublitTimeToTime()
	t2, _ := lastTime.ConvertPublitTimeToTime()
	if !t1.Equal(t2) {
		t.Error("List was not sorted properly.")
	}
}

func TestCanCreateNewStatus(t *testing.T) {
	t.Parallel()

	state := STATE_ACCEPTED
	poID := 1
	message := "My message"
	s := New(state, poID, message)

	if s.Status != state.AsString() {
		t.Errorf(`State of status did not match expected, got "%s" want "%s"`, s.Status, state.AsString())
	}

	if s.PrintOrderId != poID {
		t.Errorf(`Print order ID of status did not match expected, got "%d" want "%d"`, s.PrintOrderId, poID)
	}

	if s.Message != message {
		t.Errorf(`Message of status did not match expected, got "%s" want "%s"`, s.Message, message)
	}

	if s.SenderType != SENDER_TYPE_SUBCONTRACTOR {
		t.Errorf(`Sender type of status did not match expected, got "%s" want "%s"`, s.SenderType, SENDER_TYPE_SUBCONTRACTOR)

	}
}

func TestCanSetStatus(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint: POST}

	s := New(STATE_ACCEPTED, 1, "Some message")

	c := &MockProductionAPIClient{}

	c.PostCall = func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) {
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

func TestCanNotSetStatusForExisting(t *testing.T) {
	t.Parallel()

	s := New(STATE_ACCEPTED, 1, "Some message")
	s.ID = 1

	c := &MockProductionAPIClient{}

	err := s.Store(c)

	if err == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

func TestCanShowStatus(t *testing.T) {
	t.Parallel()
	id := 4
	r := Resource{Endpoint:SHOW,Id:id}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printorderstatus.Status" {
			t.Errorf(`Expected a Status as model but got: "%s"`, iType)
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

func TestCanIndexStatus(t *testing.T) {
	t.Parallel()
	r := Resource{Endpoint:INDEX}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printorderstatus.IndexResponse" {
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

// Test helper Client Mock
type MockProductionAPIClient struct {
	ReturnError bool
	T           *testing.T
	PostCall    func(t *testing.T, endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header))
	GetCall    func(t *testing.T, endpoint production.Endpointer, payload interface{}, queryParams ...func(q url.Values))
}

func (c *MockProductionAPIClient) Post(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	if c.ReturnError {
		return errors.New("Some error")
	}

	if c.PostCall != nil {
		c.PostCall(c.T, endpoint, payload, result, headers...)
	}

	return nil
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

func Example_indexStatusOnOrderId() {
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

	// Create filter that only fetches statuses for print order with id 5678.
	printOrderId := 5678
	filter := common.QueryAttr(
		common.AttrQuery{
			Name: PRINT_ORDER_ID,
			Value: fmt.Sprint(printOrderId),// Convert int to string
		},
	)

	// Index Status with the filter.
	po, err := Index(c, filter)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range po.Data {
		// This would output something like: State: Accepted, message: Some message, Status reported: 2017-07-17 13:31:05
		fmt.Printf("State: %s, message: %s, Status reported: %v\n", v.Status, v.Message, v.CreatedAt)
	}
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
		fmt.Printf("Could not find Status with id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find PrintOrder for id. The below would output true.
	fmt.Println(po.ID==id)
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

	// Index Status with limit
	po, err := Index(c, common.QueryLimit(10, 0))
	if err != nil {
		log.Fatal(err)
	}

	// Prints out number of returned items in response.
	fmt.Printf("Total matches: %d, number of items in list: %d\n", po.Count, len(po.Data))
}

func ExampleStatus_Store() {
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

	printorderId := 1234
	s := New(STATE_ACCEPTED, printorderId, "")

	err := s.Store(c)
	if err != nil {
		log.Fatal(err)
	}

	// If request is ok, the below would print the newly created status ID.
	fmt.Println(s.ID)
}

func ExampleStatusList_GetLast() {
	sl := StatusList{
		{ID:1, UpdatedAt:"2017-07-14 09:02:13"},
		{ID:3, UpdatedAt:"2017-07-17 22:12:10"},
		{ID:5, UpdatedAt:"2017-07-15 13:32:21"},
	}

	s := sl.GetLast()

	fmt.Println(s.ID)
	// Output: 3
}