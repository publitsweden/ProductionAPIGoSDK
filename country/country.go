// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit ProductionAPI Country resource.
//
// Package contains methods for indexing and showing countries.
package country

import (
	"fmt"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/url"
	"strings"
)

// Country attributes constants.
const (
	ID         = "id"
	Name       = "name"
	NativeName = "native_name"
	ISO2       = "iso2"
	ISO3       = "iso3"
	ISONUM     = "isonum"
)

// Country, as defined in the Publit API.
type Country struct {
	ID         int    `json:"id,string"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
	ISO2       string `json:"iso2"`
	ISO3       string `json:"iso3"`
	ISONUM     string `json:"isonum"`
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// Endpoint enumeration type.
type Endpoint int

// Endpoint enum constants.
const (
	INDEX Endpoint = 1 + iota
	SHOW
)

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX: "countries",
	SHOW:  "countries/%v",
}

// Returns Country from Publits production API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*Country, error) {
	co := &Country{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, co, queryParams...)
	return co, err
}

type CountriesList []*Country

// Index response object.
type IndexResponse struct {
	Count int           `json:"count"`
	Next  string        `json:"next"`
	Prev  string        `json:"prev"`
	Data  CountriesList `json:"data"`
}

// Indexes Countries from the Publit API.
// Returns IndexResponse where data contains StatusList.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
}

// Resource struct
type Resource struct {
	Endpoint Endpoint
	Id       int
}

// Method to Resource that fullfils the Enpointer interface as stated in production.
func (r Resource) GetEndpoint() string {
	e := endpoints[r.Endpoint]

	end := e
	if strings.Contains(e, "%v") && r.Id != 0 {
		end = fmt.Sprintf(e, r.Id)
	}

	return end
}
