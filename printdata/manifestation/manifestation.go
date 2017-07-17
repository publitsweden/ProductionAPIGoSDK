// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit manifestation data.
// Placed under printdata since it's only accessible via the printdata resource in the Publit production API.
package manifestation

import (
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata/manifestation/isbn"
)

// Holds manifestation information based on the Publit production APIs response.
type Manifestation struct {
	ID          int               `json:"id,string,omitempty"`
	WorkId      int               `json:"work_id,string"`
	ProductId   int               `json:"product_id,string"`
	IsbnID      int               `json:"isbn_id,string"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	PublishedAt common.PublitTime `json:"published_at"`
	CreatedAt   common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt   common.PublitTime `json:"updated_at,omitempty"`
	DeletedAt   common.PublitTime `json:"deleted_at"`
	Format      string            `json:"format"`
	Isbn        *isbn.ISBN        `json:"isbn"`
}
