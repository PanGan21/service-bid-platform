package testdata

import "time"

var today = time.Now()
var twoDaysAgo = today.AddDate(0, 0, -2).UTC().UnixMilli()
var yesterday = today.AddDate(0, 0, -1).UTC().UnixMilli()
var tomorrow = today.AddDate(0, 0, 1).UTC().UnixMilli()

var MockRequest = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo"}
var MockRequestYesterday = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo"}
var MockRequestTomorrow = map[string]interface{}{"Title": "mockTitle", "Postcode": "12345", "Info": "mockInfo"}

var MockRejectionReason = map[string]interface{}{"RejectionReason": "mockRejectionReason"}
