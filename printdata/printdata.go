// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit print order data print data.
// PrintData contains information about particular print items. Such as file references and sizes.
// Each PrintOrder has a list of PrintData defining quantities and what should be printed.
//
// This package contains methods to show and index print order print data.
package printdata

import (
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"github.com/publitsweden/ProductionAPIGoSDK/file"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata/bookbinding"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata/manifestation"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata/printitempaper"
	"net/url"
	"strings"
)

// With constants. Use for loading relations from PrintData.
const (
	WITH_MANIFESTATION      = "manifestation"
	WITH_MANIFESTATION_ISBN = "manifestation.isbn"
	WITH_FILE               = "file"
	WITH_PRINT_ITEM_PAPER   = "print_item_paper"
	WITH_PRINT_ITEM         = "print_item_paper.print_item"
	WITH_BOOK_BINDING       = "book_binding"
)

// Resource struct.
type Resource struct {
	Endpoint Endpoint
	Id       int
}

// Endpoint enumeration type.
type Endpoint int

// Endpoints.
const (
	INDEX Endpoint = 1 + iota
	SHOW
)

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX: "print_order_print_data",
	SHOW:  "print_order_print_data/%v",
}

// PrintDataList type.
type PrintDataList []*PrintData

// PrintData attribute constants
const (
	ID                  = "id"
	PRINT_ORDER_ID      = "print_order_id"
	MANIFESTATION_ID    = "manifestation_id"
	FILE_ID             = "file_id"
	PRINT_ITEM_PAPER_ID = "print_item_paper_id"
	BOOK_BINDING_ID     = "book_binding_id"
	AMOUNT              = "amount"
	PAGES               = "pages"
	WIDTH               = "width"
	HEIGHT              = "height"
	COLOR_PAGES_AMOUNT  = "color_pages_amount"
	COLOR_PAGES         = "color_pages"
	REFERENCE_NUMBER    = "reference_number"
	COLOR_PRINT         = "color_print"
	LENGTH_UNIT         = "length_unit"
	FORMAT              = "format"
	PUBLISHER           = "publisher"
	TITLE               = "title"
	SUBTITLE            = "subtitle"
	EDGE_WIDTH          = "edgewidth"
	CREATED_AT          = "created_at"
	UPDATED_AT          = "updated_at"
)

// Holds print data information based on the Publit production APIs "print_data" resource response.
type PrintData struct {
	ID               int                            `json:"id,string,omitempty"`
	PrintOrderID     int                            `json:"print_order_id,string"`
	ManifestationID  int                            `json:"manifestation_id,string"`
	FileID           int                            `json:"file_id,string"`
	PrintItemPaperID int                            `json:"print_item_paper_id,string"`
	BookBindingID    int                            `json:"book_binding_id,string"`
	Amount           int                            `json:"amount,string"`
	Pages            int                            `json:"pages,string"`
	Width            float64                        `json:"width"`
	Height           float64                        `json:"height"`
	ColorPagesAmount int                            `json:"color_pages_amount,string"`
	ColorPages       string                         `json:"color_pages"`
	ReferenceNumber  string                         `json:"reference_number"`
	ColorPrint       common.PublitBool              `json:"color_print"`
	LengthUnit       string                         `json:"length_unit"`
	Format           string                         `json:"format"`
	Publisher        string                         `json:"publisher"`
	Title            string                         `json:"title"`
	Subtitle         string                         `json:"subtitle"`
	EdgeWidth        float64                        `json:"edgewidth"`
	File             *file.File                     `json:"file,omitempty"`
	Manifestation    *manifestation.Manifestation   `json:"manifestation,omitempty"`
	PrintItemPaper   *printitempaper.PrintItemPaper `json:"print_item_paper,omitempty"`
	BookBinding      *bookbinding.BookBinding       `json:"book_binding,omitempty"`
	CreatedAt        common.PublitTime              `json:"created_at,omitempty"`
	UpdatedAt        common.PublitTime              `json:"updated_at,omitempty"`
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// Returns PrintDAta from Publit API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*PrintData, error) {
	pd := &PrintData{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, pd, queryParams...)
	return pd, err
}

// Index response object.
type IndexResponse struct {
	Count int           `json:"count"`
	Next  string        `json:"next"`
	Prev  string        `json:"prev"`
	Data  PrintDataList `json:"data"`
}

// Indexes PrintData from the Publit API.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
}

// Compiles PrintDataList to amp indexed on manifestation id.
func (data PrintDataList) GetPrintDataPerManifestation() map[int][]*PrintData {
	pd := make(map[int][]*PrintData)

	if len(data) > 0 {
		for _, v := range data {
			pd[v.ManifestationID] = append(pd[v.ManifestationID], v)
		}
	}

	return pd
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
