package request

import "github.com/PanGan21/pkg/entity"

type RequestEvents interface {
	PublishRequestApproved(request *entity.Request, timestamp int64) error
}
