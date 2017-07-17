package file

import (
	"bytes"
	"errors"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"github.com/publitsweden/APIUtilityGoSDK/client"
	"fmt"
	"log"
)

func TestCanIndexFiles(t *testing.T) {
	t.Parallel()
	r := Resource{Endpoint: INDEX}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "file.IndexResponse" {
			t.Errorf(`Expected a File.IndexResponse as model but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	qp := func(q url.Values) {}
	_, err := Index(c, qp)

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

func TestCanShowFile(t *testing.T) {
	t.Parallel()
	id := 5
	r := Resource{Endpoint: SHOW, Id: id}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "file.File" {
			t.Errorf(`Expected a File struct but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	qp := func(q url.Values) {}
	_, err := Show(c, id, qp)

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

func TestCanShowFileWithPresigned(t *testing.T) {
	t.Parallel()
	id := 5
	r := Resource{Endpoint: SHOW, Id: id}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		iType := reflect.Indirect(reflect.ValueOf(model)).Type().String()
		if iType != "file.File" {
			t.Errorf(`Expected a File struct but got: "%s"`, iType)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")

			qp := queryParams[0]
			q := url.Values{}
			qp(q)

			if q.Get(common.QUERY_KEY_AUX) != AUX_PRESIGNED {
				t.Error("Expected query parameter for presigned_url to be set but wasn not.")
			}
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	_, err := Show(c, id, GetPresignedAuxParamFunc())

	if err != nil {
		t.Error("Got Error but did not expect one")
	}
}

func TestCanGetPresignedURLForFile(t *testing.T) {
	t.Parallel()
	f := &File{ID: 1}

	r := Resource{Endpoint: SHOW, Id: 1}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match URL. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")

			qp := queryParams[0]
			q := url.Values{}
			qp(q)

			if q.Get(common.QUERY_KEY_AUX) != AUX_PRESIGNED {
				t.Error("Expected query parameter for presigned_url to be set but wasn not.")
			}
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	f.GetPresignedUrl(c)
}

func TestWorkerCanFetchPresignedForFileList(t *testing.T) {
	t.Parallel()
	fl := FileList{
		&File{ID: 1},
		&File{ID: 2},
		&File{ID: 3},
		&File{ID: 4},
	}

	// Set up the expected endpoints for the file list above.
	expectedEndPoints := []Resource{
		{Endpoint: SHOW, Id: 1},
		{Endpoint: SHOW, Id: 2},
		{Endpoint: SHOW, Id: 3},
		{Endpoint: SHOW, Id: 4},
	}

	// Create comma separated endpoint list for use in error text.
	endpointString := ""
	for _, v := range expectedEndPoints {
		endpointString += "," + v.GetEndpoint()
	}

	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {

		endpointMatch := false
		for _, v := range expectedEndPoints {
			if endpoint.GetEndpoint() == v.GetEndpoint() {
				endpointMatch = true
			}
		}

		if !endpointMatch {
			t.Errorf(`Endpoint did not match any expected URLs. Got "%s", expected one of "%s"`, endpoint.GetEndpoint(), endpointString)
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")

			qp := queryParams[0]
			q := url.Values{}
			qp(q)

			if q.Get(common.QUERY_KEY_AUX) != AUX_PRESIGNED {
				t.Error("Expected query parameter for presigned_url to be set but wasn not.")
			}
		}
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	errorMap := fl.GetPresigned(c)

	for _, v := range errorMap {
		if v != nil {
			t.Error("Got error but was not expecting one.", v.Error())
		}
	}
}

func TestCanDownloadFilesFromList(t *testing.T) {
	urlList := []struct {
		Url  string
		Body []byte
	}{
		{"some/already/set/presignedurl_1", []byte(`body of file 1.`)},
		{"some/already/set/presignedurl_2", []byte(`body of file 2.`)},
	}

	// Change the PlainGetter method to a mock
	PlainGetter = func(url string) (*http.Response, error) {
		resp := &http.Response{}

		body := []byte(``)
		for _, v := range urlList {
			if v.Url == url {
				body = v.Body
			}
		}

		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		resp.StatusCode = http.StatusOK

		return resp, nil
	}

	// Create fileList with files to download.
	fl := FileList{
		&File{ID: 1, OriginalName: "somefile1.txt"},
		&File{ID: 2, OriginalName: "somefile2.txt", Presigned: urlList[1].Url},
	}

	// We're only expecting endpoint for file with id=1. Since it is the only one that does not contain a presigned url.
	r := Resource{Endpoint: SHOW, Id: 1}

	// Create callback checks for the initial presigned url request.
	cb := func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) {
		if endpoint.GetEndpoint() != r.GetEndpoint() {
			t.Errorf(`Endpoint did not match any expected URLs. Got "%s", expected "%s"`, endpoint.GetEndpoint(), r.GetEndpoint())
		}

		// Assume send one query parameter
		if len(queryParams) != 1 {
			t.Error("No query params detected, but was expecting one.")

			qp := queryParams[0]
			q := url.Values{}
			qp(q)

			if q.Get(common.QUERY_KEY_AUX) != AUX_PRESIGNED {
				t.Error("Expected query parameter for presigned_url to be set but wasn not.")
			}
		}

		// Set the wanted unpresigned file here
		// Note we're assuming that it is the file with index 0 that should be set here.
		fl[0].Presigned = urlList[0].Url
	}

	c := &MockProductionAPIClient{
		ReturnError: false,
		T:           t,
		GetCall:     cb,
	}

	// Create output directory
	outdir, err := ioutil.TempDir("", "outputdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outdir)

	// Run method to test
	errorMap, err := fl.DownloadFiles(c, outdir)

	if err != nil {
		t.Error("Received an error but was not expecting one.", err.Error())
	}

	for _, v := range errorMap {
		if v != nil {
			t.Error("Got error but was not expecting one.", v.Error())
		}
	}

	// Check files are there
	dlFiles, err := ioutil.ReadDir(outdir)
	if err != nil {
		t.Fatalf("Could not read output directory: %v", outdir)
	}

	if len(dlFiles) != 2 {
		t.Errorf("Expected two files to be downloaded but got: %v",len(dlFiles))
	}

	for _, df := range dlFiles {
		path := outdir+"/"+df.Name()

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			t.Fatalf("Could not open file: %v", path)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatalf("Could not read file: %v", path)
		}

		match := false
		for _, v := range urlList {
			if string(v.Body) == string(b) {
				match = true
			}
		}

		if !match {
			t.Error("File contents did not match any of the expected.")
		}
	}
}

func TestDownloadFileReturnsErrorIfCouldNOtRunGetPresigned(t *testing.T) {
	t.Parallel()
	// Create fileList with files to download.
	fl := FileList{
		&File{ID: 1, OriginalName: "somefile1.txt"},
	}

	c := &MockProductionAPIClient{
		ReturnError: true,
	}

	// Create output directory
	outdir, err := ioutil.TempDir("", "outputdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outdir)

	// Run method to test
	errorMap, err := fl.DownloadFiles(c, outdir)

	if err != nil {
		t.Error("Received an error but was not expecting one.", err.Error())
	}

	//Since we're only providing one file in file list we know the length and index of the map.
	e := errorMap[fl[0].ID]
	if e == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

func TestDownloadFileReturnsErrorIfFileCouldNotBeDownloadedDueToBadStatus(t *testing.T) {
	// Change the PlainGetter method to a mock
	PlainGetter = func(url string) (*http.Response, error) {
		resp := &http.Response{}
		resp.StatusCode = http.StatusBadRequest

		return resp, nil
	}

	// Create fileList with files to download.
	fl := FileList{
		&File{ID: 1, OriginalName: "somefile1.txt", Presigned:"some/url/to/presigned"},
	}

	c := &MockProductionAPIClient{}

	// Create output directory
	outdir, err := ioutil.TempDir("", "outputdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outdir)

	// Run method to test
	errorMap, err := fl.DownloadFiles(c, outdir)

	if err != nil {
		t.Error("Received an error but was not expecting one.", err.Error())
	}

	//Since we're only providing one file in file list we know the length and index of the map.
	e := errorMap[fl[0].ID]
	if e == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

func TestDownloadFilesReturnsErrorIfDirectoryDoesNotExists(t *testing.T) {
	t.Parallel()
	unexistingDir := "some/dir/that/doesnt/exist"
	c := &MockProductionAPIClient{}

	fl := FileList{&File{}}

	_, err := fl.DownloadFiles(c, unexistingDir)

	if err == nil {
		t.Error("Did not receive an error but was expecting one.")
	}
}

// Test helper Client Mock
type MockProductionAPIClient struct {
	ReturnError bool
	T           *testing.T
	GetCall     func(t *testing.T, endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values))
}

func (c *MockProductionAPIClient) Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error {
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
	f, err := Show(c, id)
	if err != nil {
		fmt.Printf("Could not find id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find file the below should output true.
	fmt.Println(f.ID==id)
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

	// Create a limit to fetch a maximum of 10 items
	limit := common.QueryLimit(10, 0)

	// Run Index to retrieve list of Files.
	f, err := Index(c, limit)

	if err != nil {
		fmt.Printf("Could not list files. Got errors: %v\n", err.Error())
	}

	// Print the number of DeliveryNumbers in list from response.
	fmt.Println(len(f.Data))
}

func ExampleFile_GetPresignedUrl() {
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

	// Create an "existing" file by setting ID.
	id := 1
	f := &File{ID: id}

	// Get Presigned URL.
	err := f.GetPresignedUrl(c)
	if err != nil {
		fmt.Printf("Could not find id: %d. Got errors: %v\n", id, err.Error())
	}

	// If request to Publit could find file the presigned_url of the file should now be available.
	fmt.Println(f.Presigned)
}

func ExampleFileList_GetPresigned() {
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

	// Filter request to show only created files numbers after 2017-06-01 00:00:00.
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

	// Run Index to retrieve list of Files.
	f, err := Index(c, filter, limit)
	if err != nil {
		log.Fatal("Could not retrieve list of files.", err.Error())
	}

	// Retrieve FileList from Index response.
	fl := f.Data

	// Get presigned for the list of files.
	errs := fl.GetPresigned(c)

	// Range the error map to see if all files could receive their presigned urls.
	for k, v := range errs {
		if v != nil {
			fmt.Printf("Got error when trying to retrieve fie list for file id: %d, error: %v\n",k, v)
		}
	}
}

func ExampleFileList_DownloadFiles() {
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

	// Create FileList.
	// See the example for FileList.GetPresigned() for a more exhaustive example.
	fl := FileList{
		{ID: 1},
		{ID: 2},
	}

	outputPath := "some/ptah/to/folder"

	// All of the FileList operations use workers to concurrently handle FileList data.
	// The default is set to 5 worker but can be altered with the SetWorkerAmount method.
	SetWorkerAmount(3)

	// Get presigned for the list of files.
	errs, err := fl.DownloadFiles(c, outputPath)

	// If method itself returned an error the abort all further execution (or return early).
	if err != nil {
		log.Fatal(err)
	}

	// Range the error map to see if all files could be downloaded.
	for k, v := range errs {
		if v != nil {
			fmt.Printf("Got error when trying to download file with id: %d, error: %v\n",k, v)
		}
	}

	// Check files are there
	dlFiles, err := ioutil.ReadDir(outputPath)
	if err != nil {
		log.Fatalf("Could not read output directory: %v", outputPath)
	}

	// print amount of files in output path.
	fmt.Println(len(dlFiles))
}