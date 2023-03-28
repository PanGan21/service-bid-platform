package request

import "github.com/PanGan21/pkg/entity"

type RequestEvents interface {
	PublishRequestCreated(request *entity.Request) error
}
