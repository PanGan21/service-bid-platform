package testdata

import "github.com/google/uuid"

var DefaultRoles []int
var MockUser = map[string]interface{}{"Username": uuid.New().String(), "Email": uuid.New().String(), "Phone": uuid.New().String(), "Password": "mockPassword"}
var AdminUser = map[string]interface{}{"Username": "SuperAdmin", "Password": "password"}
