package printdata

import (
	"testing"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/url"
	"errors"
	"reflect"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/client"
)

func TestCanGroupPrintDataOnManifestation(t *testing.T) {
	t.Parallel()
	pd := PrintDataList{
		{
			ID: 1,
			ManifestationID: 1,
		},
		{
			ID: 2,
			ManifestationID: 1,
		},
		{
			ID: 3,
			ManifestationID: 2,
		},
		{
			ID: 4,
			ManifestationID: 2,
		},
	}

	grouped := pd.GetPrintDataPerManifestation()

	if len(grouped) != 2 {
		t.Errorf("Expected array length of 2, but got %v.",len(grouped))
	}

	for k, v := range grouped {
		for _, p := range v {
			if p.ManifestationID != k {
				t.Errorf("Grouped print data had non expected manifestation id. Have %v want %v.",p.ManifestationID, k)
			}
		}
	}
}

func TestCanIndexPrintData(t *testing.T) {
	t.Parallel()

	r := Resource{Endpoint:INDEX}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printdata.IndexResponse" {
			t.Errorf(`Expected a PrintData index response as model but got: "%s"`, iType)
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
		t.Error("Received an error but was not expecting one.")
	}
}

func TestCanShowPrintData(t *testing.T) {
	t.Parallel()

	id := 7
	r := Resource{Endpoint:SHOW,Id:id}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "printdata.PrintData" {
			t.Errorf(`Expected a PrintData as model but got: "%s"`, iType)
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
	_, err := Show(c, id, qp)

	if err != nil {
		t.Error("Received an error but was not expecting one.")
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

	orderId := 1
	// Filter request to show print data for order with id 1.
	filter := common.QueryAttr(
		common.AttrQuery{
			Name: PRINT_ORDER_ID,
			Value: fmt.Sprint(orderId), //Convert int to string
		},
	)

	// Run Index to retrieve list of PrintData.
	pd, err := Index(c, filter)

	if err != nil {
		fmt.Printf("Could not list print order print data after 2017-01-01. Got errors: %v\n", err.Error())
	}

	// Print the number of DeliveryNumbers in list from response.
	fmt.Println(len(pd.Data))
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

	PrintDataId := 1

	// Run Show to retrieve PrintData with id 1.
	pd, err := Show(c, PrintDataId)

	if err != nil {
		fmt.Printf("Could not find print order print data with id 1. Got errors: %v\n", err.Error())
	}

	// Prints true for succesful response.
	fmt.Println(pd.ID==PrintDataId)
}

func ExamplePrintDataList_GetPrintDataPerManifestation() {
	pdl := PrintDataList{
		{ManifestationID: 1, ID: 1},
		{ManifestationID: 1, ID: 2},
		{ManifestationID: 2, ID: 3},
		{ManifestationID: 2, ID: 4},
	}

	pdmap := pdl.GetPrintDataPerManifestation()

	for k, v := range pdmap {
		for _, pd := range v {
			fmt.Printf("manifestation_id: %d, has PrintData with ID: %d\n",k, pd.ID)
		}
	}

	// Output:
	// manifestation_id: 1, has PrintData with ID: 1
	// manifestation_id: 1, has PrintData with ID: 2
	// manifestation_id: 2, has PrintData with ID: 3
	// manifestation_id: 2, has PrintData with ID: 4
}