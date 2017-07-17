// Copyright 2017 Publit Sweden AB. All rights reserved.

// Manages Publit print order delivery numbers.
// Package contains methods to show, index, store, update and delete delivery numbers.
package deliverynumber

import (
	"errors"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/http"
	"net/url"
	"strings"
)

// DeliveryNumber attribute constants.
const (
	ID              = "id"
	PRINT_ORDER_ID  = "print_order_id"
	DELIVERY_NUMBER = "delivery_number"
	MESSAGE         = "message"
	CREATED_AT      = "created_at"
	UPDATED_AT      = "updated_at"
)

// Print order delivery number struct.
type DeliveryNumber struct {
	ID             int               `json:"id,string,omitempty"`
	PrintOrderID   int               `json:"print_order_id,string"`
	DeliveryNumber string            `json:"delivery_number"`
	Message        string            `json:"message,omitempty"`
	CreatedAt      common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt      common.PublitTime `json:"updated_at,omitempty"`
}

// Resource struct
type Resource struct {
	Endpoint Endpoint
	Id       int
}

// Endpoint enumeration type.
type Endpoint int

// Endpoint enum constants.
const (
	INDEX Endpoint = 1 + iota
	SHOW
	POST
	PUT
	DELETE
)

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX:  "print_order_delivery_numbers",
	SHOW:   "print_order_delivery_numbers/%v",
	POST:   "print_order_delivery_numbers",
	PUT:    "print_order_delivery_numbers/%v",
	DELETE: "print_order_delivery_numbers/%v",
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// ProductionAPIPoster defines how the client should perform POST calls.
type ProductionAPIPoster interface {
	Post(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error
}

// ProductionAPIPoster defines how the client should perform PUT calls.
type ProductionAPIPutter interface {
	Put(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error
}

// ProductionAPIDeleter defines how the client should perform DELETE calls
type ProductionAPIDeleter interface {
	Delete(endpoint production.Endpointer, result interface{}, headers ...func(h *http.Header)) error
}

// Creates new DeliveryNumber and returns pointer.
// Mainly used for storing "new" DeliveryNumbers.
func New(printOrderId int, deliveryNumber, message string) *DeliveryNumber {
	d := &DeliveryNumber{
		PrintOrderID:   printOrderId,
		DeliveryNumber: deliveryNumber,
	}

	if message != "" {
		d.Message = message
	}

	return d
}

// Returns DeliveryNumber from Publit API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*DeliveryNumber, error) {
	po := &DeliveryNumber{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, po, queryParams...)
	return po, err
}

// Delivery number list type.
type DeliveryNumberList []*DeliveryNumber

// Index response object.
type IndexResponse struct {
	Count int                `json:"count"`
	Next  string             `json:"next"`
	Prev  string             `json:"prev"`
	Data  DeliveryNumberList `json:"data"`
}

// Indexes DeliveryNumbers from the Publit API.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
}

// Updates DeliveryNumber.
func (d *DeliveryNumber) Update(c ProductionAPIPutter) error {
	if d.ID == 0 {
		return errors.New("Can not update a non existing number. (ID is missing).")
	}

	r := Resource{Endpoint: PUT, Id: d.ID}
	return c.Put(r, d, d)
}

// Stores delivery number.
func (d *DeliveryNumber) Store(c ProductionAPIPoster) error {
	if d.ID != 0 {
		return errors.New("Can not create new delivery number for an existing one. (ID is set).")
	}
	r := Resource{Endpoint: POST}
	dnl := &DeliveryNumberList{d}
	return c.Post(r, d, dnl)
}

// Deletes delivery number.
func (d *DeliveryNumber) Delete(c ProductionAPIDeleter) error {
	if d.ID == 0 {
		return errors.New("Can not DELETE a non existing number. (ID is missing).")
	}
	r := Resource{Endpoint: DELETE, Id: d.ID}
	return c.Delete(r, d)
}

// Method to Resource that fulfils the Endpointer interface as stated in production.
func (r Resource) GetEndpoint() string {
	e := endpoints[r.Endpoint]

	end := e
	if strings.Contains(e, "%v") && r.Id != 0 {
		end = fmt.Sprintf(e, r.Id)
	}

	return end
}
