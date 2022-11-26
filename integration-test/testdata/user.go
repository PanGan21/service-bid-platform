package testdata

import "github.com/google/uuid"

var defaultRoles = make([]int, 0)
var MockUser = map[string]interface{}{"username": uuid.New(), "password": uuid.New()}
