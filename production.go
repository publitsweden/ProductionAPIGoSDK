// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit ProductionAPI client and calls.
//
// See the subdirectories (packages) for more information on real world usage.
package production

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"net/http"
	"net/url"
)

const (
	API                  = "production"
	API_VERSION          = "v2.0"
	RESOURCE_STATUSCHECK = "status_check"
)

type Endpointer interface {
	GetEndpoint() string
}

// APICaller is an interface that defines how a client should use the Publit APIs.
// The github.com/publitsweden/APIUtilityGoSDK/client.Client fulfills this interface.
type APICaller interface {
	Call(r *http.Request) (*http.Response, error)
	CallRaw(r *http.Request) (*http.Response, error)
	SetNewAPIToken(r *http.Request) error
}

// APIClient hold Client information for connecting to the Publit APIs and base URLs.
type APIClient struct {
	Client  APICaller
	BaseUrl string
}

// StatusCheck checks if the Publit service is up.
func (c *APIClient) StatusCheck() bool {
	url := c.compileStatusCheckURL()

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return false
	}

	// Use CallRaw since no authentication is needed for status check.
	r, err := c.Client.CallRaw(req)

	if err != nil {
		return false
	}

	if r.StatusCode != http.StatusOK {
		return false
	}

	return true
}

// Compiles statuscheck URL against the production API.
func (c APIClient) compileStatusCheckURL() string {
	return fmt.Sprintf("%s/%s/%s", c.BaseUrl, API_VERSION, RESOURCE_STATUSCHECK)
}

// Performs a GET method action against the Publit production API.
func (c APIClient) Get(endpoint Endpointer, model interface{}, queryParams ...func(q url.Values)) error {
	endUrl := c.CompileEndpointURL(endpoint.GetEndpoint())
	req, _ := http.NewRequest(http.MethodGet, endUrl, nil)

	q := req.URL.Query()
	for _, v := range queryParams {
		v(q)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Call(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return MakeResponseError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(model)

	if err != nil {
		return err
	}

	return nil
}

// Performs a POST method action against the Publit production API.
func (c APIClient) Post(endpoint Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	return c.postPut(http.MethodPost, endpoint, payload, result, headers...)
}

// Performs a PUT method action against the Publit production API.
func (c APIClient) Put(endpoint Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	return c.postPut(http.MethodPut, endpoint, payload, result, headers...)
}

// Performs a post or put method action against the Publit production API.
func (c APIClient) postPut(method string, endpoint Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error {
	endUrl := c.CompileEndpointURL(endpoint.GetEndpoint())

	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, _ := http.NewRequest(method, endUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	h := &req.Header
	for _, v := range headers {
		v(h)
	}

	resp, err := c.Client.Call(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return MakeResponseError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		return err
	}

	return nil
}

// Performs a DELETE http call against the Publit production API.
func (c APIClient) Delete(endpoint Endpointer, result interface{}, headers ...func(h *http.Header)) error {
	endUrl := c.CompileEndpointURL(endpoint.GetEndpoint())
	req, _ := http.NewRequest(http.MethodDelete, endUrl, nil)

	h := &req.Header
	for _, v := range headers {
		v(h)
	}

	resp, err := c.Client.Call(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MakeResponseError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		return err
	}

	return nil
}

// Compiles regular endpoints URL.
func (c APIClient) CompileEndpointURL(endpoint string) string {
	return fmt.Sprintf("%v/%v/%v/%v", c.BaseUrl, API, API_VERSION, endpoint)
}

// Attempts to make a better response error from response.
func MakeResponseError(resp *http.Response) error {
	if resp.Header.Get("Content-Type") == "application/json" {
		APIErr := &common.APIErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(APIErr)
		if err == nil && APIErr.HasInformation() { // Only return this error message if APIErr has information.
			return APIErr.GetAsError()
		}
	}
	return errors.New(fmt.Sprintf(`Response not ok. No information given. Code: "%v"`, resp.StatusCode))
}
