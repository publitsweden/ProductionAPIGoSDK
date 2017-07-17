// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit print order data.
// The package contains method to show and index print orders.
package printorder

import (
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"github.com/publitsweden/ProductionAPIGoSDK/country"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata"
	"github.com/publitsweden/ProductionAPIGoSDK/printorderstatus"
	"net/url"
	"strings"
)

// With constants.
// Defines available relations to load for GET requests.
// Helps with making filters, but can also be customized.
const (
	WITH_STATUSES                      = "print_order_statuses"
	WITH_PRINT_DATA                    = "print_order_print_data"
	WITH_PRINT_DATA_FILE               = "print_order_print_data.file"
	WITH_PRINT_DATA_MANIFESTATION      = "print_order_print_data.manifestation"
	WITH_PRINT_DATA_MANIFESTATION_ISBN = "print_order_print_data.manifestation.isbn"
	WITH_PRINT_DATA_PRINT_ITEM_PAPER   = "print_order_print_data.print_item_paper"
	WITH_PRINT_DATA_PRINT_ITEM         = "print_order_print_data.print_item_paper.print_item"
	WITH_PRINT_DATA_BOOK_BINDING       = "print_order_print_data.book_binding"
	WITH_DELIVERY_COUNTRY              = "delivery_country"
)

// Endpoint enumeration constants
const (
	INDEX Endpoint = 1 + iota
	SHOW
)

// Resource struct
type Resource struct {
	Endpoint Endpoint
	Id       int
}

// Endpoint enumeration type.
type Endpoint int

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX: "print_orders",
	SHOW:  "print_orders/%v",
}

// PrintOrder attribute constants.
const (
	ID                     = "id"
	INTERMEDIATOR_REF      = "intermediator_order_reference"
	CLIENT_REF             = "client_order_reference"
	DELIVERY_MSG           = "delivery_message"
	ORDER_WEIGHT           = "order_weight"
	BULKY                  = "is_bulky"
	RECIPIENT_FIRSTNAME    = "firstname"
	RECIPIENT_LASTNAME     = "lastname"
	RECIPIENT_COMPANY_NAME = "delivery_company_name"
	DELIVERY_STREET        = "delivery_street"
	DELIVERY_ZIP           = "delivery_zip"
	DELIVERY_CITY          = "delivery_city"
	DELIVERY_PHONE         = "delivery_phone_number"
	DELIVERY_COUNTRY_ID    = "delivery_country_id"
	DELIVERY_PRE_PAID      = "delivery_address_pre_paid"
	STATUS                 = "status"
	ACTIVE                 = "active"
	EXPECTED_SHIP_DATE     = "expected_shipment_date"
	CREATED_AT             = "created_at"
	UPDATED_AT             = "updated_at"
)

// PrintOrder struct. Holds information available in the Publit print_orders endpoint.
type PrintOrder struct {
	ID                   int                         `json:"id,string,omitempty"`
	IntermediatorRef     string                      `json:"intermediator_order_reference"`
	ClientRef            string                      `json:"client_order_reference"`
	DeliveryMsg          string                      `json:"delivery_message"`
	OrderWeight          string                      `json:"order_weight"`
	Bulky                common.PublitBool           `json:"is_bulky"`
	RecipientFirstname   string                      `json:"firstname"`
	RecipientLastname    string                      `json:"lastname"`
	RecipientCompanyName string                      `json:"delivery_company_name"`
	DeliveryStreet       string                      `json:"delivery_street"`
	DeliveryZip          string                      `json:"delivery_zip"`
	DeliveryCity         string                      `json:"delivery_city"`
	DeliveryPhone        string                      `json:"delivery_phone_number"`
	DeliveryCountryId    string                      `json:"delivery_country_id"`
	DeliveryPrePaid      string                      `json:"delivery_address_pre_paid"`
	Status               string                      `json:"status"`
	Active               common.PublitBool           `json:"active"`
	ExpectedShipDate     common.PublitTime           `json:"expected_shipment_date"`
	CreatedAt            common.PublitTime           `json:"created_at,omitempty"`
	UpdatedAt            common.PublitTime           `json:"updated_at,omitempty"`
	Statuses             printorderstatus.StatusList `json:"print_order_statuses,omitempty"`
	PrintData            printdata.PrintDataList     `json:"print_order_print_data,omitempty"`
	DeliveryCountry      *country.Country            `json:"delivery_country,omitempty"`
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// Returns PrintOrder from Publit API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*PrintOrder, error) {
	po := &PrintOrder{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, po, queryParams...)
	return po, err
}

// Index response object.
type IndexResponse struct {
	Count int          `json:"count"`
	Next  string       `json:"next"`
	Prev  string       `json:"prev"`
	Data  []PrintOrder `json:"data"`
}

// Indexes PrintOrders from the Publit API.
// Returns IndexResponse where data contains PrintOrder list.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
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
