package testdata

import "github.com/google/uuid"

var DefaultRoles []int
var MockUser = map[string]interface{}{"Username": uuid.New().String(), "Password": "mockPassword"}
