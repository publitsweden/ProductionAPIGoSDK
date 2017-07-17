// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit ISBN data.
// Placed under printdata/manifestation since it's only accessible via the printdata.manifestation relation in the Publit production API.
package isbn

import "github.com/publitsweden/APIUtilityGoSDK/common"

// Holds ISBN information based on the Publit production APIs response
type ISBN struct {
	ID            int               `json:"id,string,omitempty"`
	FormattedISBN string            `json:"formatted_isbn"`
	AccountId     int               `json:"account_id,string"`
	ContractorId  int               `json:"contractor_id,string"`
	CreatedAt     common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt     common.PublitTime `json:"updated_at,omitempty"`
}
