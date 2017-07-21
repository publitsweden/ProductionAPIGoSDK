package production_test

import (
	. "github.com/publitsweden/ProductionAPIGoSDK"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/APILog"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK/file"
	"github.com/publitsweden/ProductionAPIGoSDK/printorder"
	"github.com/publitsweden/ProductionAPIGoSDK/printorderstatus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestCanCheckStatus(t *testing.T) {
	t.Parallel()
	t.Run(
		"If status is ok",
		func(t *testing.T) {
			caller := &MockAPICaller{}

			caller.Response = createCallerResponse(http.StatusOK, "")

			baseurl := "somebaseurl"

			c := &APIClient{caller, baseurl}

			ok := c.StatusCheck()

			if !ok {
				t.Error("Expected status check to pass, but received false.")
			}
		},
	)

	t.Run(
		"If status is not ok",
		func(t *testing.T) {
			caller := &MockAPICaller{}

			caller.Response = createCallerResponse(http.StatusBadRequest, "")

			baseurl := "somebaseurl"

			c := &APIClient{caller, baseurl}

			ok := c.StatusCheck()

			if ok {
				t.Error("Expected status check to fail, but received true.")
			}
		},
	)

	t.Run(
		"If call returns error",
		func(t *testing.T) {
			caller := &MockAPICaller{}

			caller.Response = createCallerResponse(http.StatusBadRequest, "")
			caller.ReturnErrors = true

			baseurl := "somebaseurl"

			c := &APIClient{caller, baseurl}

			ok := c.StatusCheck()

			if ok {
				t.Error("Expected status check to fail, but received true.")
			}
		},
	)

}

func TestCanPerformGetRequest(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}

	caller.Response = createCallerResponse(http.StatusOK, `{"some":"body"}`)

	baseurl := "somebaseurl"

	c := &APIClient{caller, baseurl}

	int := &struct {
		Some string `json:"some"`
	}{}

	err := c.Get(NewEndpoint(), int)

	if err != nil {
		t.Error("Expected Get to pass but received error.", err.Error())
	}

	if int.Some != "body" {
		t.Error("Unmarshalled struct did not match expected.")
	}
}

func TestGetReturnsErrorIfCallFails(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}

	caller.Response = createCallerResponse(http.StatusOK, `{"some":"body"}`)

	caller.ReturnErrors = true

	baseurl := "somebaseurl"

	c := &APIClient{caller, baseurl}

	int := &struct{}{}
	err := c.Get(NewEndpoint(), int)

	if err == nil {
		t.Error("Expected an error due to call failed but did not receive one.")
	}
}

func TestGetReturnsErrorIfStatusCodeNotOk(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}

	caller.Response = createCallerResponse(http.StatusBadRequest, `{"some":"body"}`)

	baseurl := "somebaseurl"

	c := &APIClient{caller, baseurl}

	int := &struct{}{}
	err := c.Get(NewEndpoint(), int)

	if err == nil {
		t.Error("Expected an error due to status not ok but did not receive one.")
	}
}

func TestGetReturnsErrorIfBodyCanNotBeUnmarshalled(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}

	caller.Response = createCallerResponse(http.StatusOK, `{"some","much:faulty":,%€%;"body"}`)

	baseurl := "somebaseurl"

	c := &APIClient{caller, baseurl}

	int := &struct{}{}
	err := c.Get(NewEndpoint(), int)

	if err == nil {
		t.Error("Expected an error due to unmarshalling errors but did not receive one.")
	}
}

func TestCanPerformPOSTRequest(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}
	caller.T = t

	i := struct {
		Name string `json:"name"`
	}{Name: "test"}
	j := &i

	caller.CallTestCallback = func(t *testing.T, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Got error was not expecting one.")
		}

		ic := i

		json.Unmarshal(b, ic)

		if ic.Name != i.Name {
			t.Error("Request body did not match expected.")
		}
	}

	baseurl := "somebaseurl"
	caller.Response = createCallerResponse(http.StatusOK, `{"name":"newTestName"}`)

	c := &APIClient{caller, baseurl}

	err := c.Post(NewEndpoint(), &i, j)

	if err != nil {
		t.Error("Received an error but was not expecting to.")
	}

	if i.Name != "newTestName" {
		t.Error("Struct did not have expected value.")
	}
}

func TestCanPerformPUTRequest(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}
	caller.T = t

	i := struct {
		Name string `json:"name"`
	}{Name: "test"}
	j := &i

	caller.CallTestCallback = func(t *testing.T, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Got error was not expecting one.")
		}

		ic := i

		json.Unmarshal(b, ic)

		if ic.Name != i.Name {
			t.Error("Request body did not match expected.")
		}
	}

	baseurl := "somebaseurl"
	caller.Response = createCallerResponse(http.StatusOK, `{"name":"newTestName"}`)

	c := &APIClient{caller, baseurl}

	err := c.Put(NewEndpoint(), &i, j)

	if err != nil {
		t.Error("Received an error but was not expecting to.")
	}

	if i.Name != "newTestName" {
		t.Error("Struct did not have expected value.")
	}
}

func TestCanPerformDeleteRequest(t *testing.T) {
	t.Parallel()
	caller := &MockAPICaller{}
	caller.T = t

	i := struct {
		Name string `json:"name"`
	}{}

	baseurl := "somebaseurl"
	caller.Response = createCallerResponse(http.StatusOK, `{"name":"newTestName"}`)

	c := &APIClient{caller, baseurl}

	err := c.Delete(NewEndpoint(), &i)

	if err != nil {
		t.Error("Received an error but was not expecting to.")
	}

	if i.Name != "newTestName" {
		t.Error("Struct did not have expected value.")
	}
}

func TestPostPutErrors(t *testing.T) {
	t.Parallel()
	table := []struct {
		TestName string
		TestFunc func(t *testing.T)
	}{
		{
			TestName: "If model can not be marshalled into json",
			TestFunc: func(t *testing.T) {
				caller := &MockAPICaller{}

				// Can not serialize and marshal a chan.
				i := make(chan int)

				baseurl := "somebaseurl"

				c := &APIClient{caller, baseurl}

				// Run method through POST (would be just as fine with PUT
				err := c.Post(NewEndpoint(), &i, &i)
				if err == nil {
					t.Error("Did not receive an error, but was expecting to.")
				}
			},
		},
		{
			TestName: "If Call returns error",
			TestFunc: func(t *testing.T) {
				caller := &MockAPICaller{}
				caller.ReturnErrors = true
				caller.Response = createCallerResponse(http.StatusOK, "")

				baseurl := "somebaseurl"
				c := &APIClient{caller, baseurl}

				i := struct{}{}
				// Run method through POST (would be just as fine with PUT
				err := c.Post(NewEndpoint(), &i, &i)

				if err == nil {
					t.Error("Did not receive an error, but was expecting to.")
				}
			},
		},
		{
			TestName: "If response status code is not ok",
			TestFunc: func(t *testing.T) {
				caller := &MockAPICaller{}
				caller.Response = createCallerResponse(http.StatusBadRequest, "")

				baseurl := "somebaseurl"
				c := &APIClient{caller, baseurl}

				i := struct{ Name string }{}
				err := c.Post(NewEndpoint(), &i, &i)

				if err == nil {
					t.Error("Did not receive an error, but was expecting to.")
				}
			},
		},
		{
			TestName: "If response json can not be marshalled to interface",
			TestFunc: func(t *testing.T) {
				caller := &MockAPICaller{}
				caller.Response = createCallerResponse(http.StatusOK, `{"somwire,ddgd:""jsonstructur,,,:"newTestName"}`)

				baseurl := "somebaseurl"
				c := &APIClient{caller, baseurl}

				i := struct {
					Name string `json:"name"`
				}{Name: "test"}
				err := c.Post(NewEndpoint(), &i, &i)

				if err == nil {
					t.Error("Did not receive an error, but was expecting to.")
				}
			},
		},
	}

	for _, v := range table {
		t.Run(
			v.TestName,
			v.TestFunc,
		)
	}
}

func TestCanMakeResponseError(t *testing.T) {
	t.Parallel()

	t.Run(
		"By parsing response",
		func(t *testing.T) {
			errorMessage := []byte(`{"Code":400,"Type":"BadRequest","Errors":[{"Info":"Some error","Type":"BadRequest"}],"CombinedInfo":"Some error"}`)

			resp := &http.Response{
				Status:     "400 BadRequest",
				StatusCode: http.StatusBadRequest,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: ioutil.NopCloser(bytes.NewBuffer(errorMessage)),
			}

			e := MakeResponseError(resp)

			if e == nil {
				t.Error("No error recieved but was expecting one.")
			}

			expected := `Code: "400", Type: "BadRequest", Combined info: "Some error"`
			if e.Error() != expected {
				t.Error("Error message did not match expected.")
			}
		},
	)
	t.Run(
		"Unauthorized error without body",
		func(t *testing.T) {
			errorMessage := []byte(`{}`)

			resp := &http.Response{
				Status:     "401 Unauthorized",
				StatusCode: http.StatusUnauthorized,
				Body: ioutil.NopCloser(bytes.NewBuffer(errorMessage)),
			}

			e := MakeResponseError(resp)

			if e == nil {
				t.Error("No error recieved but was expecting one.")
			}

			expected := `Unauthorized. Code: "401"`
			if e.Error() != expected {
				t.Errorf(`Error message did not match expected. Got: "%v", Expected "%v"`, e.Error(), expected)
			}
		},
	)
	t.Run(
		"Default error without body and not unauthorized.",
		func(t *testing.T) {
			errorMessage := []byte(`{}`)

			resp := &http.Response{
				Status:     "400 BadRequest",
				StatusCode: http.StatusBadRequest,
				Body: ioutil.NopCloser(bytes.NewBuffer(errorMessage)),
			}

			e := MakeResponseError(resp)

			if e == nil {
				t.Error("No error recieved but was expecting one.")
			}

			expected := `Response not ok. No information given. Code: "400"`
			if e.Error() != expected {
				t.Errorf(`Error message did not match expected. Got: "%v", Expected "%v"`, e.Error(), expected)
			}
		},
	)
}

func createCallerResponse(status int, body string) *http.Response {
	resp := &http.Response{}
	resp.StatusCode = status

	if body != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	}

	return resp
}

type MockAPICaller struct {
	ReturnErrors     bool
	Response         *http.Response
	CallTestCallback func(t *testing.T, r *http.Request)
	T                *testing.T
}

func (c *MockAPICaller) Call(r *http.Request) (*http.Response, error) {
	if c.ReturnErrors {
		return c.Response, errors.New("Some error")
	}

	if c.CallTestCallback != nil {
		c.CallTestCallback(c.T, r)
	}

	return c.Response, nil
}

func (c *MockAPICaller) CallRaw(r *http.Request) (*http.Response, error) {
	return c.Call(r)
}

func (c *MockAPICaller) SetNewAPIToken(r *http.Request) error {
	if c.ReturnErrors {
		return errors.New("Some error")
	}

	return nil
}

// Creates new endpoint.
func NewEndpoint() Endpoint { return 1 }

// For fulfilling the endpointer interface.
type Endpoint int

// For fulfilling the endpointer interface.
func (e Endpoint) GetEndpoint() string {
	return "someendpoint"
}

// Examples

// Check if service is up.
func Example_checkStatus() {
	//Create client.
	c := APIClient{
		Client: client.New(
			func(c *client.Client) {
				c.User = "myusername"
				c.Password = "mypassword"
			},
		),
		BaseUrl: "https://url.to.publit",
	}

	// Check if service is up.
	ok := c.StatusCheck()

	if !ok {
		log.Fatal("Service is not up.")
	}

	log.Println("Service is up!")

	// Do something
}

// Ingest PrintOrder, set status and download files for print.
func Example_ingestOrder() {
	// Set logging output and output level.
	APILog.LogOutput = os.Stdout
	APILog.OutputLevel = APILog.LEVEL_DEBUG

	// Create client.
	c := APIClient{
		Client: client.New(
			func(c *client.Client) {
				c.User = "myusername"
				c.Password = "mypassword"
			},
		),
		BaseUrl: "https://url.to.publit",
	}

	// Check if service is up.
	// This will also set the client token which will be used for authenticating further requests.
	ok := c.StatusCheck()

	if !ok {
		log.Fatal("Service is not up.")
	}

	// Get print order by id. And load it together with the relation to PrintData.File and PrintItemPaper.PrintItem and DeliveryCountry.
	// The PrintData.File relation is to later download the items.
	// The PrintItemPaper.PrintItem is for determining what papers the print items should be printed on.
	// The DeliveryCountry is for knowing where the parcel should be sent once the print is complete.
	printorderID := 1234
	po, err := printorder.Show(c, printorderID, common.QueryWith(
		printorder.WITH_PRINT_DATA_FILE,
		printorder.WITH_PRINT_DATA_PRINT_ITEM,
		printorder.WITH_PRINT_DATA_BOOK_BINDING,
		printorder.WITH_DELIVERY_COUNTRY))
	if err != nil {
		log.Fatal(err)
	}

	// Range PrintData to get the individual files.
	var fl file.FileList
	for _, v := range po.PrintData {
		fl = append(fl, v.File)
	}

	// Download files to folder.
	outputPath := "some/path/to/download/folder"
	errmap, err := fl.DownloadFiles(c, outputPath)
	if err != nil {
		log.Fatal(err)
	}

	// Crude error handling.
	// Aborts if any errors were found and removes the output folder.
	for _, v := range errmap {
		if v != nil {
			os.RemoveAll(outputPath)
			log.Fatal(v)
		}
	}

	// Set status to "accepted" once order has been ingested.
	s := printorderstatus.New(printorderstatus.STATE_ACCEPTED, printorderID, "Ingested print order.")
	err = s.Store(c)
	if err != nil {
		log.Fatal(err)
	}
}

// Converts Publit standard error response to an error.
func ExampleMakeResponseError() {
	resp := &http.Response{
		Status:     "400 BadRequest",
		StatusCode: http.StatusBadRequest,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"Code":400,"Type":"BadRequest","Errors":[{"Info":"Some error","Type":"BadRequest"}],"CombinedInfo":"Some error"}`))),
	}

	err := MakeResponseError(resp)

	fmt.Println(err.Error())
	// Output: Code: "400", Type: "BadRequest", Combined info: "Some error"
}

// Compiles endpoint URL based on baseurl and endpoint.
func ExampleAPIClient_CompileEndpointURL() {
	c := APIClient{
		Client: client.New(
			func(c *client.Client) {
				c.User = "myusername"
				c.Password = "mypassword"
			},
		),
		BaseUrl: "https://url.to.publit",
	}

	endpoint := "someendpoint"

	url := c.CompileEndpointURL(endpoint)

	fmt.Println(url)
	// Output: https://url.to.publit/production/v2.0/someendpoint
}
