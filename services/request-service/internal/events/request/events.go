package request

import "github.com/PanGan21/pkg/entity"

type RequestEvents interface {
	PublishRequestCreated(topic string, request *entity.Request) error
}
