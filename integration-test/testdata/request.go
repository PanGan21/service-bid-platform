package testdata

import "time"

var today = time.Now()
var yesterday = today.AddDate(0, 0, -1).Unix()
var tomorrow = today.AddDate(0, 0, 1).Unix()

var MockRequest = map[string]interface{}{"Title": "mockTitle", "Postcode": "mockPostcode", "Info": "mockInfo", "Deadline": 1}
var MockRequestYesterday = map[string]interface{}{"Title": "mockTitle", "Postcode": "mockPostcode", "Info": "mockInfo", "Deadline": yesterday}
var MockRequestTomorrow = map[string]interface{}{"Title": "mockTitle", "Postcode": "mockPostcode", "Info": "mockInfo", "Deadline": tomorrow}
