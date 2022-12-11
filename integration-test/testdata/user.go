package testdata

import "github.com/google/uuid"

var DefaultRoles []int
var MockUser = map[string]interface{}{"username": uuid.New().String(), "password": "mockPassword"}
