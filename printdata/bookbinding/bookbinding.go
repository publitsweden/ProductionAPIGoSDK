// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit book binding data.
// Placed under printdata since it's only accessible via the printdata resource in the Publit production API.
package bookbinding

import (
	"github.com/publitsweden/APIUtilityGoSDK/common"
)

// Holds PrintItemPaper information based on the Publit production APIs response.
type BookBinding struct {
	ID        int               `json:"id,string,omitempty"`
	Type      string            `json:"type,string"`
	CreatedAt common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt common.PublitTime `json:"updated_at,omitempty"`
}
