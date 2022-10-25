package testdata

import "github.com/google/uuid"

var MockUser = map[string]interface{}{"username": uuid.New(), "password": uuid.New()}
