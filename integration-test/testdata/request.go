package testdata

import "github.com/google/uuid"

var MockRequest = map[string]interface{}{"title": uuid.New(), "postcode": uuid.New(), "info": uuid.New(), "deadline": 1}
