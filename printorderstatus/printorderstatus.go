// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit print order data statuses.
//
// Print order statuses is used in Publit to communicate the state of the print order.
// Whenever a print order changes status it should be communicated to Publit.
//
// This package contains methods to show, index and store print order statuses.
package printorderstatus

import (
	"errors"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// General constants for status.
const (
	// The type of sender for the Production API. Only valid sender_type to send via the Production API is "Subcontractor".
	SENDER_TYPE_SUBCONTRACTOR = "Subcontractor"
)

// PrintOrder attribute constants.
const (
	ID             = "id"
	PRINT_ORDER_ID = "print_order_id"
	SENDER_TYPE    = "sender_type"
	STATUS         = "status"
	MESSAGE        = "message"
	CREATED_AT     = "created_at"
	UPDATED_AT     = "updated_at"
)

// Print order status, as defined in the Publit Production API.
type Status struct {
	ID           int               `json:"id,string,omitempty"`
	PrintOrderId int               `json:"print_order_id,string"`
	SenderType   string            `json:"sender_type"`
	Status       string            `json:"status"`
	Message      string            `json:"message,omitempty"`
	CreatedAt    common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt    common.PublitTime `json:"updated_at,omitempty"`
}

// State enumeration type.
type State int

// State constants.
const (
	STATE_EXPORTED State = 1 + iota
	STATE_ACCEPTED
	STATE_IN_PRODUCTION
	STATE_SENT
	STATE_PARTIALLY_SENT
	STATE_DELIVERED
	STATE_ABORTED
	STATE_RETURNED
	STATE_RESEND
)

// Statuses map.
var statues map[State]string = map[State]string{
	STATE_EXPORTED:       "Exported",
	STATE_ACCEPTED:       "Accepted",
	STATE_IN_PRODUCTION:  "In production",
	STATE_SENT:           "Sent",
	STATE_PARTIALLY_SENT: "Partially sent",
	STATE_DELIVERED:      "Delivered",
	STATE_ABORTED:        "Aborted",
	STATE_RETURNED:       "Returned",
	STATE_RESEND:         "Resend",
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// ProductionAPIPoster defines how the client should perform POST calls.
type ProductionAPIPoster interface {
	Post(endpoint production.Endpointer, payload interface{}, result interface{}, headers ...func(h *http.Header)) error
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
)

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX: "print_order_statuses",
	SHOW:  "print_order_statuses/%v",
	POST:  "print_order_statuses",
}

// Creates new status.
func New(state State, PrintOrderId int, message string) *Status {
	s := &Status{
		Status:       state.AsString(),
		PrintOrderId: PrintOrderId,
		SenderType:   SENDER_TYPE_SUBCONTRACTOR,
	}

	if message != "" {
		s.Message = message
	}

	return s
}

// Status list type.
type StatusList []*Status

// Sets status.
func (s *Status) Store(c ProductionAPIPoster) error {
	if s.ID != 0 {
		return errors.New("Can not create new status for an existing one. (ID is set).")
	}
	r := Resource{Endpoint: POST}
	pr := &StatusList{s}
	err := c.Post(r, s, pr)
	return err
}

// Returns Status from Publits production API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*Status, error) {
	po := &Status{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, po, queryParams...)
	return po, err
}

// Index response object.
type IndexResponse struct {
	Count int        `json:"count"`
	Next  string     `json:"next"`
	Prev  string     `json:"prev"`
	Data  StatusList `json:"data"`
}

// Indexes Statuses from the Publit API.
// Returns IndexResponse where data contains StatusList.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
}

// Returns state as readable string as defined by the Publit APIs.
func (s State) AsString() string {
	return statues[s]
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

// Retrieves last reported status from StatusList.
func (l StatusList) GetLast() *Status {
	sort.Sort(l)
	return l[len(l)-1]
}

// Len method to fulfil sort.Sort interface.
func (p StatusList) Len() int {
	return len(p)
}

// Less method to fulfil sort.Sort interface.
func (p StatusList) Less(i, j int) bool {
	t, _ := p[i].UpdatedAt.ConvertPublitTimeToTime()
	t2, _ := p[j].UpdatedAt.ConvertPublitTimeToTime()
	return t.Before(t2)
}

// Swap method to fulfil sort.Sort interface.
func (p StatusList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
