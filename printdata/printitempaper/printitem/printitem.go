// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit print item data.
// Placed under printdata/printitempaper since it's only accessible via the printdata resource in the Publit production API.
package printitem

// Holds PrintItem information based on the Publit production APIs response.
type PrintItem struct {
	ID   int    `json:"id,string,omitempty"`
	Type string `json:"type,string"`
}
