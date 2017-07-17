// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit print item paper data.
// Placed under printdata since it's only accessible via the printdata resource in the Publit production API.
package printitempaper

import (
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK/printdata/printitempaper/printitem"
)

// Holds PrintItemPaper information based on the Publit production APIs response.
type PrintItemPaper struct {
	ID                   int                  `json:"id,string,omitempty"`
	PrintItemID          int                  `json:"print_item_id,string"`
	Name                 string               `json:"name"`
	ProprietaryPaperName string               `json:"proprietary_paper_name,string"`
	PaperCode            string               `json:"paper_code"`
	Bulk                 string               `json:"bulk"`
	Weight               string               `json:"weight"`
	Description          string               `json:"description"`
	PrintItem            *printitem.PrintItem `json:"print_item,omitempty"`
	CreatedAt            common.PublitTime    `json:"created_at,omitempty"`
	UpdatedAt            common.PublitTime    `json:"updated_at,omitempty"`
}
