package testdata

import "time"

var today = time.Now()
var twoDaysAgo = today.AddDate(0, 0, -2).UTC().UnixMilli()
var yesterday = today.AddDate(0, 0, -1).UTC().UnixMilli()
var tomorrow = today.AddDate(0, 0, 1).UTC().UnixMilli()

var MockAuction = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo", "Deadline": twoDaysAgo}
var MockAuctionYesterday = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo", "Deadline": yesterday}
var MockAuctionTomorrow = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo", "Deadline": tomorrow}

var MockRejectionReason = map[string]interface{}{"RejectionReason": "mockRejectionReason"}
